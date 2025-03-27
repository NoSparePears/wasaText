<template>
    <div class="group-view">
        <header v-if="name">
            <img :src="avatar" alt="Group photo" class="group-photo" />
            <h1>{{ name }}</h1>
            <button @click="goToInfo" class="button">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#info"></use></svg>
              
            </button>
        </header>
        <div class="messages">
            <p v-if="!messages.length">No messages yet</p>
            <ul>
                <li v-for="msg in messages" :key="msg.message.msgID"
                    :class="{'own-message': isOwnMessage(msg.message.senderID), 'other-message': !isOwnMessage(msg.message.senderID)}" @click="openMessageOptions(msg.message)">
                    
                    <!-- Display sender's name -->
                    <span class="message-sender">{{ msg.user.username }}</span>
                    <div class="message-wrapper">
                        <div class="message-bubble">
                            <span class="message-content">{{ msg.message.content }}</span>
                        </div> 
                        <span class="message-timestamp">{{ formatTimestamp(msg.message.timestamp) }}</span>
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





    </div>
</template>

<script>
import Search from '@/components/Search.vue';

export default {
    components: {
        Search
    },
    data() {
        return {
            name: this.$route.query.name || "Unknown",
            groupID: this.$route.query.groupID,
            avatar: this.$route.query.avatar || "default_propic.jpg",
            messages: [], // Lista dei messaggi
            text: null, // Testo del messaggio da inviare
            photo: null, // Foto da inviare
            showModal: false, // Mostra il modal per le opzioni del messaggio
            selectedMessage: null, // Messaggio selezionato per le opzioni
            errormsg: '', // Messaggio di errore
            searchModalVisible: false,  //mostra modale ricerca utenti
        }
    },
    methods: {
        async sendMessage() {
            this.errormsg = '';
            const userID = sessionStorage.getItem('id');
            const token = sessionStorage.getItem('token');
            try {
                let response = await this.$axios.post(`/profiles/${userID}/groups/${this.groupID}/messages`, { content: this.text }, {
                headers: { 'Authorization': token }
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
        async fetchMessages() {
            this.errormsg = '';
            const userID = sessionStorage.getItem('id');
            const token = sessionStorage.getItem('token');
            try {
                let response = await this.$axios.get(`/profiles/${userID}/groups/${this.groupID}/messages`, {
            headers: { 'Authorization': sessionStorage.getItem('token') }
            });
            this.messages = response.data;
            if (!this.messages) this.messages = [];
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error fetching messages';
                console.error('Error fetching messages:', error);
            }
        },
        async deleteMessage() {
            const userID = sessionStorage.getItem('id');
            const token = sessionStorage.getItem('token');
            try {
                await this.$axios.delete(`/profiles/${userID}/groups/${this.groupID}/messages/${this.selectedMessage.msgID}`, {
                    headers: { 'Authorization': token }
            });
                this.closeModal();
                this.fetchMessages();
            } catch (e) {
                console.error(e);
            }
        },
        isOwnMessage(senderID) {
            const userID = sessionStorage.getItem('id');
            return String(senderID) === String(userID);
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
        openMessageOptions(message) {
            this.selectedMessage = message;
            this.showModal = true;
        },
        closeModal() {
            this.showModal = false;
            this.selectedMessage = null;
        },
        toggleSearchModal() {
            this.searchModalVisible = !this.searchModalVisible;
        },
        goToInfo() {
            console.log('Navigating to group info');
            try {
                this.$router.push({
                    path: `/home/groups/${this.groupID}/info`,
                    query: {
                        name: this.name || "Unknown",
                        groupID: this.groupID,
                        avatar: this.avatar || "default_propic.jpg",
                    }
                })
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error opening group s info';
                console.error('Error opening group info:', error);
            }
        },
    },
    mounted() {
      // Se l'utente non è loggato, reindirizza alla pagina di login
      if (!sessionStorage.getItem('token')) {
        this.$router.push("/");
        return;
      }
      this.fetchMessages();
      this.intervalId = setInterval(async () => {
        clearInterval(this.intervalId);
        await this.fetchMessages();
        this.intervalId = setInterval(this.fetchMessages, 1000);
      }, 1000);
      
    },
    beforeUnmount() { 
        // Pulisci intervallo quando il componente viene distrutto
        if (this.intervalId) {
            clearInterval(this.intervalId);
        }
    }
}
</script>

<style scoped>
  .group-view {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    background-color: #f7f7f7;
  }
  
  .group-photo {
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
  
  .message-sender {
    font-weight: bold;
    color: #555;
    margin-bottom: 4px;
    display: block;
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
    font-size: 12px;
    color: gray;
    margin-top: 3px;
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
  </style>
