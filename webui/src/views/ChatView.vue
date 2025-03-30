<template>
    <div class="chat-view">
      <header v-if="name">
        <img :src="avatar" alt="User Avatar" class="avatar" />
        <h2>{{ name }}</h2>
      </header>
      <div class="messages">
        <p v-if="!messages.length">No messages yet.</p>
        <ul>
          <li v-for="message in messages" :key="message.message.msgID" 
              :class="{'own-message': isOwnMessage(message.message.senderID), 'other-message': !isOwnMessage(message.message.senderID)}" @click="openMessageOptions(message.message)">
            <div class="message-wrapper">
              <div class="message-bubble">
                <span class="message-content">{{ message.message.content }}</span>
              </div> 

              <div class="message-meta" v-if="isOwnMessage(message.message.senderID)">
                <span class="message-status">{{ messageStatus(message.message) }}</span>
                <span class="message-timestamp">{{ formatTimestamp(message.message.timestamp) }}</span>
              </div>
              <span v-else class="message-timestamp">{{ formatTimestamp(message.message.timestamp) }}</span>

              <!-- Display Comments -->
              <div v-if="message.comments && message.comments.length" class="message-comments">
                <div 
                  class="comment-bubble" 
                  v-for="comment in message.comments" 
                  :key="comment.commID"
                  :class="{'own-comment': isOwnComment(comment.senderID)}">
                  <strong>{{ comment.sendUsername }}:</strong> {{ comment.emoji }}
                
                  <!-- Delete button for own comments using SVG icon -->
                  <button v-if="isOwnComment(comment.senderID)" @click.stop="deleteComment(comment)" class="delete-button">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash"></use></svg>
                  </button>
                </div>
              </div>

            </div>
          </li>
        </ul>
      </div>
      <!-- Input per inviare un messaggio testuale -->
      <div class="input-group">
          <input type="text" class="form-control" v-model="text" placeholder="Type your message here">
          <button class="btn btn-outline-primary" @click="sendMessage()">Send</button>
      </div>
      <!-- Modal for message options -->
      <div v-if="showModal" class="message-options-modal">
        <div class="modal-content">
          <h3>What would you like to do with this message?</h3>
          <button @click="toggleCommentModal" class="btn btn-primary">Comment</button>
          <button @click="toggleSearchModal" class="btn btn-primary">Forward</button>
          <button @click="deleteMessage" class="btn btn-danger">Delete</button>
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
        </div>
      </div>
      <!-- Modal for searching users to whom forward a message -->
      <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="forwardMessage" title="search">
        <template v-slot:header>
          <h3>Users</h3>
        </template>
      </Search>
      <!-- Modal for commenting a message -->
      <Comment :show="commentModalVisible" :msg="selectedMessage" @close="toggleCommentModal" @emoji-selected="commentMessage" />
    </div>
    
  </template>
  
  <script>
  import Search from '@/components/Search.vue';
  import Comment from '@/components/Comment.vue';

  export default {
    components: {
      Search,
      Comment
    },
    props: ['id'],
    data() {
      return {
        name: this.$route.query.name || "Unknown",
        destID: this.$route.query.destID,
        avatar: this.$route.query.avatar || "default_propic.jpg",
        messages: [], // Lista dei messaggi
        text: null, // Testo del messaggio da inviare
        photo: null, // Foto da inviare
        showModal: false, // Mostra il modal per le opzioni del messaggio
        selectedMessage: null, // Messaggio selezionato per le opzioni
        errormsg: '', // Messaggio di errore
        searchModalVisible: false,  //mostra modale ricerca utenti
        commentModalVisible: false, //mostra modale commento
      };
    },
    methods: {
      isOwnMessage(senderID) {
        const userID = sessionStorage.getItem('id');
        return String(senderID) === String(userID);
      },
      isOwnComment(senderID) {
        const userID = sessionStorage.getItem('id');
        return String(senderID) === String(userID);
      },
      messageStatus(message) {
        if (message.received) {
          return "✓✓"; // Double check for read messages
        } else if (message.sent) {
          return "✓"; // Single check for sent messages
        }
        return ""; // No checkmark if not sent
      },
      async getMessages() {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        try {
          let response = await this.$axios.get(`/profiles/${userID}/conversations/${this.destID}/messages`, {
            headers: { 'Authorization': sessionStorage.getItem('token') }
          });
          this.messages = response.data;
          if (!this.messages) this.messages = [];
          
        } catch (error) {
          this.errormsg = error.response?.data?.message || 'Error fetching messages';
          console.error('Error fetching messages:', error);
        }
      },

      formatTimestamp(timestamp) {
        if (!timestamp) return "";

        const date = new Date(timestamp);
        const now = new Date();

        // Controlla se il messaggio è di oggi
        const isToday = date.toDateString() === now.toDateString();
        
        // Controlla se il messaggio è di ieri
        const yesterday = new Date();
        yesterday.setDate(now.getDate() - 1);
        const isYesterday = date.toDateString() === yesterday.toDateString();

        // Formatta orario (HH:MM)
        const timeString = date.toLocaleTimeString("it-IT", {hour: "2-digit", minute: "2-digit" });

        if (isToday) {
          return `Oggi, ${timeString}`;
        } else if (isYesterday) {
          return `Ieri, ${timeString}`;
        } else {
          // Formatta data (GG/MM/AAAA)
          const dateString = date.toLocaleDateString("it-IT");
          return `${dateString}, ${timeString}`;
        }

      },

      async sendMessage() {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        // aggiungi possibilità per mandare foto
        try {
            let response = await this.$axios.post(`/profiles/${userID}/conversations/${this.destID}/messages`, { content: this.text }, {
            headers: { 'Authorization': sessionStorage.getItem('token') }
            })
            // Reset the variables used for sending the message
            this.text = null;
            this.photo = null;
            // Assuming response.data is the message object itself now, not wrapped in a 'message' field
            this.messages.push(response.data);
            
        } catch (error) {
            this.errormsg = error;
            console.error('Error sending message:', error);
        }
      },

      openMessageOptions(message) {
        // Store the selected message and show the modal
        this.selectedMessage = message;
        this.showModal = true;
      },

      closeModal() {
        // Close the modal and reset the selected message
        this.showModal = false;
        this.selectedMessage = null;
      },

      toggleSearchModal() {
        this.searchModalVisible = !this.searchModalVisible;
      },
      toggleCommentModal() {
        this.commentModalVisible = !this.commentModalVisible;
      },
      async commentMessage(emoji) {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        try {
          let response = await this.$axios.put(`/profiles/${userID}/conversations/${this.selectedMessage.convoID}/messages/${this.selectedMessage.msgID}/comments`, {emoji: emoji}, {
            headers: { 'Authorization': token }
          });
          if (response.status === 200) {
            this.selectedMessage.emoji = emoji;
            this.toggleCommentModal();
            this.closeModal();
          }
          
        } catch (error) {
          this.errormsg = error.response?.data?.message || 'Error commenting message';
          console.error('Error commenting message:', error);
        }
      },
      async deleteComment(comment) {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        try {
          let response = await this.$axios.delete(`/profiles/${userID}/conversations/${this.destID}/messages/${comment.msgID}/comments/${comment.commID}`, {
            headers: { 'Authorization': token }
          });
          
          
        } catch (error) {
          this.errormsg = error.response?.data?.message || 'Error deleting comment';
          console.error('Error deleting comment:', error);
        }
      },
      /*
      async getMessageComments(msgID) {
        for (let message of this.messages) {
          try {
            const userID = sessionStorage.getItem('id'); 
            const token = sessionStorage.getItem('token'); 
          
            let response = await this.$axios.get(
              `/profiles/${userID}/conversations/${this.destID}/messages/${message.msgID}/comments`,
              { headers: { 'Authorization': token } }
            );
          
            if (response.status === 200) {
              // Convert API response to match our frontend structure
              const formattedComments = response.data.map(commentObj => ({
                username: commentObj.sendUser.username,
                emoji: commentObj.comment.emoji,
                commentID: commentObj.comment.commentID
              }));
              
            }
          } catch (error) {
            console.error(`Error fetching comments for message ${message.msgID}:`, error);
          }
        }
      },
      */
      async forwardMessage(destUser) {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');

        try {
          const openChat = await this.$axios.get(`/profiles/${userID}/conversations/${destUser.id}`, {
            headers: { 'Authorization': token }
          });
          // Check if request was successful
          if (openChat.status === 200) {
            console.log('Chat opened successfully');
          } else {
            throw new Error('Request failed with status: ${response.status}');
          }     
          const response = await this.$axios.post(`/profiles/${userID}/conversations/${destUser.id}/messages`, {content: "forwarded: "+this.selectedMessage.content}, {
            headers: { 'Authorization': token }
          });
          this.toggleSearchModal();
          this.closeModal();
        } catch (error) {
          console.error('Error in requests:', error);
        }
        
      },

      async deleteMessage() {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        try {
          let response = await this.$axios.delete(`/profiles/${userID}/conversations/${this.destID}/messages/${this.selectedMessage.msgID}`, {
            headers: { 'Authorization': token }
          });
          this.messages = [];
          if (this.messages) this.messages = response.data; 
          
          this.closeModal();
        } catch (error) {
          this.errormsg = error.response?.data?.message || 'Error deleting message';
          console.error('Error deleting message:', error);
        }
      },
      

    },
    mounted() {
      // Se l'utente non è loggato, reindirizza alla pagina di login
      if (!sessionStorage.getItem('token')) {
        this.$router.push("/");
        return;
      }
      this.getMessages();
      this.intervalId = setInterval(async () => {
        clearInterval(this.intervalId);
        await this.getMessages();
        this.intervalId = setInterval(this.getMessages, 1000);
      }, 1000);
      
    },
    beforeUnmount() { 
        // Pulisci intervallo quando il componente viene distrutto
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
    }
  };

  </script>
  
  <style scoped>
  .chat-view {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    background-color: #f7f7f7;
  }
  
  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin-right: 10px;
  }
  
  .messages ul {
    list-style-type: none;
    padding: 0;
    margin: 0;
  }
  
  .messages li {
    display: flex;
    flex-direction: column;
    margin-bottom: 10px;
  }
  
  .own-message {
    align-items: flex-end;
  }
  
  .other-message {
    align-items: flex-start;
  }
  
  .message-wrapper {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    max-width: 70%;
  }

  .message-bubble {
    padding: 10px;
    border-radius: 10px;
    background-color: #f1f0f0;
    word-wrap: break-word;
    display: inline-block;
  }

  .message-timestamp {
    font-size: 0.75rem;
    color: rgba(0, 0, 0, 0.6);
    margin-top: 2px;
    text-align: right;
    padding-right: 5px;
  }
  
  .own-message .message-bubble {
    background-color: #007bff;
    color: white;
    align-self: flex-end;
  }
  
  .other-message .message-bubble {
    background-color: #f1f1f1;
    color: black;
    text-align: left;
    align-self: flex-start;
  }

  .own-message .message-wrapper {
    align-items: flex-end;
  }

  .other-message .message-wrapper {
    align-items: flex-start;
  }
  
  .message-meta {
    display: flex;
    align-items: center;
    font-size: 0.75rem;
    color: rgba(0, 0, 0, 0.6);
    margin-top: 2px;
  }

  .message-status {
    margin-right: 5px;
    color: #007bff;
    font-weight: bold;
  }

  .input-group {
    margin-top: 20px;
    display: flex;
  }
  
  .input-group input {
    flex-grow: 1;
    padding: 10px;
  }
  
  .input-group button {
    padding: 10px;
    margin-left: 10px;
  }
  
  .message-options-modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10;
  }
  
  .modal-content {
    background-color: white;
    padding: 20px;
    border-radius: 10px;
    text-align: center;
    width: 300px; /* Make the modal smaller */
  }
  
  .modal-content button {
    margin: 10px;
  }
  
  .modal-content button.btn-primary {
    background-color: #0084ff;
    color: white;
  }
  
  .modal-content button.btn-danger {
    background-color: #f44336;
    color: white;
  }
  
  .modal-content button.btn-secondary {
    background-color: #9e9e9e;
    color: white;
  }
  .message-comments {
    margin-top: 5px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .comment-bubble {
   background-color: #f0f0f0; /* Light gray */
   padding: 8px;
   margin-top: 4px;
   border-radius: 8px;
   font-size: 14px;
   display: inline-block;
  }
  .delete-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  margin-left: 8px;
  display: flex;
  align-items: center;
}

.delete-button svg {
  width: 16px;
  height: 16px;
  transition: fill 0.2s ease-in-out;
}

.delete-button:hover svg {
  fill: #cc0000;
}
  </style>
