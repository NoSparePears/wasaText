<template>
    <div class="chat-view">
      <header v-if="name">
        <img :src="avatar" alt="User Avatar" class="avatar" />
        <h2>{{ name }}</h2>
      </header>
      <div class="messages">
        <p v-if="!messages.length">No messages yet.</p>
        <ul>
          <li v-for="message in messages" :key="message.msgID" 
              :class="{'own-message': isOwnMessage(message.senderID), 'other-message': !isOwnMessage(message.senderID)}" @click="openMessageOptions(message)">
            <div class="message-bubble">
            {{ message.content }}
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
          <button @click="forwardMessage" class="btn btn-primary">Forward</button>
          <button @click="deleteMessage" class="btn btn-danger">Delete</button>
          <button @click="closeModal" class="btn btn-secondary">Cancel</button>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    props: ['id'],
    data() {
      return {
        name: this.$route.query.name || "Unknown",
        destID: this.$route.query.destID,
        avatar: this.$route.query.avatar || "default-avatar.png",
        messages: [], // Lista dei messaggi
        text: null, // Testo del messaggio da inviare
        photo: null, // Foto da inviare
        showModal: false, // Mostra il modal per le opzioni del messaggio
        selectedMessage: null, // Messaggio selezionato per le opzioni
        errormsg: '', // Messaggio di errore
      };
    },
    methods: {
      isOwnMessage(senderID) {
        const userID = sessionStorage.getItem('id');
        return String(senderID) === String(userID);
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

      async deleteMessage() {
        this.errormsg = '';
        const userID = sessionStorage.getItem('id');
        const token = sessionStorage.getItem('token');
        try {
          let response = await this.$axios.delete(`/profiles/${userID}/conversations/${this.destID}/messages/${this.selectedMessage.msgID}`, {
            headers: { 'Authorization': sessionStorage.getItem('token') }
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
    align-items: flex-start;
    margin-bottom: 10px;
  }
  
  .own-message {
    align-items: flex-end;
  }
  
  .other-message {
    align-items: flex-start;
  }
  
  .message-bubble {
    padding: 10px 15px;
    border-radius: 15px;
    max-width: 60%;
    word-wrap: break-word;
  }
  
  .own-message .message-bubble {
    background-color: #007bff;
    color: white;
    text-align: right;
    align-self: flex-end;
  }
  
  .other-message .message-bubble {
    background-color: #f1f1f1;
    color: black;
    text-align: left;
    align-self: flex-start;
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
