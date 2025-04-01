<template>
    <div class="chat-view">
      <header v-if="name">
        <section class="profile-picture">
          <img v-if="avatar" :src="avatar" alt="User Avatar" class="pfp-img" />
        </section>
        <h2>{{ name }}</h2>
      </header>
      <div class="messages">
        <p v-if="!messages.length">No messages yet.</p>
        <ul>
          <li v-for="message in messages" :key="message.message.msgID" 
              :class="{'own-message': isOwnMessage(message.message.senderID), 'other-message': !isOwnMessage(message.message.senderID)}" 
              @click="openMessageOptions(message.message)">
            <div class="message-wrapper">
              <div class="message-bubble">
                <!-- Forwarded Label -->
                <span v-if="message.message.isForwarded" class="forwarded-label">Forwarded</span>

                <!-- Check if message contains a photo -->
                <div v-if="message.message.isPhoto">
                  <img :src="`data:image/jpeg;base64,${message.message.content}`" alt="Sent Image" class="message-image" />
                </div>
                <span v-else class="message-content">{{ message.message.content }}</span>
                
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
      <!-- Input per inviare un messaggio -->
      <div class="input-group">
        <!-- Text Input for Messages -->
        <input type="text" class="form-control" v-model="text" placeholder="Type your message here">

        <!-- Hidden File Input for Image Upload -->
        <input type="file" ref="fileInput" @change="handleFileUpload" accept="image/*" style="display: none;">

        <!-- Plus Icon to Trigger File Upload -->
        <button class="btn btn-outline-secondary upload-btn" @click="triggerFileUpload">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle"></use></svg>
        </button>
      
        <!-- Send Button -->
        <button class="btn btn-outline-primary" @click="sendMessage">Send</button>
      </div>

      <!-- Image Preview (Only shown if an image is selected) -->
      <div v-if="photoPreview" class="image-preview">
        <img :src="photoPreview" alt="Selected Image" class="preview-img">
        <button class="btn btn-danger btn-sm remove-btn" @click="clearImage">X</button>
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
        avatar: '',
        messages: [], // Lista dei messaggi
        text: null, // Testo del messaggio da inviare
        photo: null, // Foto da inviare
        photoPreview: null,
        showModal: false, // Mostra il modal per le opzioni del messaggio
        selectedMessage: null, // Messaggio selezionato per le opzioni
        errormsg: '', // Messaggio di errore
        searchModalVisible: false,  //mostra modale ricerca utenti
        commentModalVisible: false, //mostra modale commento
      };
    },
    methods: {
      async getProfilePicture() {
        const userID = this.destID;
        const token = sessionStorage.getItem('token');
        try {
          // Fai una chiamata API per ottenere la foto del profilo
          const response = await this.$axios.get(`/profiles/${userID}/photo`, 
            { headers: { 'Authorization': token } 
          });
          // Extract base64 data from response
          if (response.data && response.data.profile_picture) {
            this.avatar = `data:image/jpeg;base64,${response.data.profile_picture}`;
          } else {
            console.error("Profile picture data is missing in the response");
          }
        } catch (error) {
          console.error("Errore nel recupero della foto del profilo:", error);
        }
      },
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
        if (!this.text && !this.photo) {
          return; // Prevent sending empty messages
        }
        try {
          let formData = new FormData();

          if (this.photo) {
            formData.append('image', this.photo);
            formData.append("isPhoto", 1);
          }
          if (this.text) {
            formData.append('content', this.text);
            formData.append("isPhoto", 0);
          }
          let response = await this.$axios.post(
            `/profiles/${userID}/conversations/${this.destID}/messages`,
            formData,
            {
              headers: {
                'Authorization': token,
                'Content-Type': 'multipart/form-data',
              },
            }
          );
          // Reset the variables used for sending the message
          this.text = null;
          this.photo = null;
          this.$refs.fileInput.value = ""; // Clear file input
          this.clearImage(); // Clear image preview
          // Assuming response.data is the message object itself now, not wrapped in a 'message' field
          this.messages.push(response.data);
          this.getMessages(); // Refresh messages
        } catch (error) {
            this.errormsg = error;
            console.error('Error sending message:', error);
        }
      },
      // Open file selection when the plus icon is clicked
      triggerFileUpload() {
        this.$refs.fileInput.click();
      },
    
      // Handle the selected image and show a preview
      handleFileUpload(event) {
        const file = event.target.files[0];
        if (file) {
          this.photo = file;
          this.photoPreview = URL.createObjectURL(file); // Create a preview URL
        }
      },
    
      // Clear selected image
      clearImage() {
        this.photo = null;
        this.photoPreview = null;
        this.$refs.fileInput.value = "";
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
          const response = await this.$axios.post(`/profiles/${userID}/conversations/${this.destID}/messages/${this.selectedMessage.msgID}`, {destID: destUser.id}, {
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
      this.getProfilePicture();
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
  
  .pfp-img {
  width: 120%; /* Ensures full coverage */
  height: 120%;
  object-fit: cover; /* Ensures image fits without stretching */
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
  .forwarded-label {
    font-size: 12px;
    color: #666;
    font-weight: bold;
    display: block;
    margin-bottom: 5px;
}


.message-content {
    font-size: 14px;
}

.message-image {
    max-width: 100%;
    border-radius: 5px;
    display: block;
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
    background-color: #f1f1f1;
    padding: 5px;
    border-radius: 5px;
    margin-bottom: 5px;
    display: flex;
    align-items: center;
}

.own-comment {
    background-color: #dcf8c6;
}
.feather {
    width: 16px;
    height: 16px;
    stroke: red;
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
.upload-btn svg {
  width: 24px;
  height: 24px;
}

.image-preview {
  margin-top: 10px;
  display: flex;
  align-items: center;
}

.preview-img {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 5px;
  margin-right: 10px;
}

.remove-btn {
  border: none;
  cursor: pointer;
}
  </style>
