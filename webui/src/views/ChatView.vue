<template>
    <div class="chat-view">
      <header v-if="name">
        <img :src="avatar" alt="User Avatar" class="avatar" />
        <h2>{{ name }}</h2>
      </header>
      <div class="messages">
        <p v-if="!messages.length">No messages yet.</p>
        <ul>
          <li v-for="message in messages" :key="message.msgID">
            {{ message.content }}
          </li>
        </ul>
      </div>
      <!-- Input per inviare un messaggio testuale -->
      <div class="input-group">
          <input type="text" class="form-control" v-model="text" placeholder="Type your message here">
          <button class="btn btn-outline-primary" @click="sendMessage()">Send</button>
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
        messages: [],
        // Testo del messaggio da inviare
        text: null,
        // Foto del messaggio da inviare
        photo: null,
      };
    },
    methods: {
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
      }

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
  
  <style>
  .avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
  }
  </style>
  