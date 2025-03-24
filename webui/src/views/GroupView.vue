<template>
    <div class="group-view">
        <header> v-if="name">
            <img :src="avatar" alt="Group photo" class="group-photo" />
            <h1>{{ name }}</h1>
        </header>
        <div class="messaages">
            <p v-if="!messages.length">No messages yet</p>
            <ul>
                <li v-for="message in messages" :key="message.id"
                    :class="{'own-message': isOwnMessage(message.senderID), 'other-message': !isOwnMessage(message.senderID)}" @click="openMessageOptions(message)">
                    <div class="message-wrapper">
                        <div class="message-bubble">
                            <span class="message-content">{{ message.content }}</span>
                        </div> 
                        <span class="message-timestamp">{{ formatTimestamp(message.timestamp) }}</span>
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
            destID: this.$route.query.destID,
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
            try {
                await this.$axios.post(`/groups/${this.$route.params.id}/messages`, {
                    content: this.text
                });
                this.text = "";
                this.fetchMessages();
            } catch (e) {
                console.error(e);
            }
        },
        async fetchMessages() {
            try {
                let response = await this.$axios.get(`/groups/${this.$route.params.id}/messages`);
                this.messages = response.data;
            } catch (e) {
                console.error(e);
            }
        },
        async deleteMessage() {
            try {
                await this.$axios.delete(`/groups/${this.$route.params.id}/messages/${this.selectedMessage.id}`);
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
        }
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
        this.intervalId = setInterval(this.getMessages, 1000);
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