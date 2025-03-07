<template>
    <div class="home-view">
      <header>
        <h1>Chat</h1>
        <div class="header-buttons">
        <button @click="fetchConversations">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"></use></svg>
          Refresh
        </button>
        <button @click="showNewChatModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Chat
        </button>
    </div>
    </header>
    <div class="filter-buttons">
      <button :class="{ active: filter === 'all' }" @click="setFilter('all')">Tutte</button>
      <button :class="{ active: filter === 'groups' }" @click="setFilter('groups')">Gruppi</button>
    </div>
    <section class="conversations-list">
      <h2>Tutte le Conversazioni</h2>
      <ul>
        <li v-for="conversation in filteredConversations" :key="conversation.id">
          <span>{{ conversation.name }}</span>
        </li>
      </ul>
    </section>
    <div v-if="isModalVisible" class="modal">
      <div class="modal-content">
        <span class="close" @click="closeNewChatModal">&times;</span>
        <h2>Nuova Chat</h2>
        <form @submit.prevent="createChat">
          <input v-model="newChatName" placeholder="Nome Chat" required />
          <select v-model="chatType" required>
            <option value="conversation">Conversazione</option>
            <option value="group">Gruppo</option>
          </select>
          <button type="submit">Crea</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      
      conversations: [],
      newChatName: '',
      chatType: 'conversation',
      isModalVisible: false,
      filter: 'all'
    };
  },
  computed: {
    filteredConversations() {
      if (this.filter === 'all') {
        return this.conversations;
      } else if (this.filter === 'groups') {
        return this.conversations.filter(conversation => conversation.type === 'group');
      }
      return [];
    }
  },
  methods: {
    //ottiene le conversazione dell utente
    async fetchConversations() {
      try {
        //effettuo richiesta api
        let response = await this.$axios.get(`/profiles/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
        //salvo dati conversazioni
        this.conversations = response.data
      } catch (error) {
        console.error('Errore nel recupero delle conversazioni:', error);
      }
    },
    async createChat() {
      try {
        const endpoint = this.chatType === 'conversation' ? '/api/conversations' : '/api/groups';
        const response = await fetch(endpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ name: this.newChatName })
        });
        if (response.ok) {
          this.fetchConversations();
          this.newChatName = '';
          this.closeNewChatModal();
        } else {
          console.error('Errore nella creazione della chat');
        }
      } catch (error) {
        console.error('Errore nella creazione della chat:', error);
      }
    },
    showNewChatModal() {
      this.isModalVisible = true;
    },
    closeNewChatModal() {
      this.isModalVisible = false;
    },
    setFilter(filter) {
      this.filter = filter;
    }
  },
  mounted() {
    // Se l'utente non è loggato, reindirizza alla pagina di login
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }
    this.fetchConversations();
    this.intervalId = setInterval(async () => {
      clearInterval(this.intervalId);
      await this.getConversations();
      this.intervalId = setInterval(this.getConversations, 1000);
    }, 1000);
  },
  beforeUnmount() {
    // Pulisci l'intervallo quando il componente viene distrutto
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }
};
</script>

<style scoped>
.home-view {
  padding: 20px;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-buttons {
  display: flex;
  gap: 10px;
}
.header-buttons button {
  background: none;
  border: none;
  cursor: pointer;
}
.header-buttons .icon {
  width: 32px;
  height: 32px;
}
.filter-buttons {
  margin-top: 10px;
  display: flex;
  gap: 10px;
}
.filter-buttons button {
  padding: 10px 20px;
  border: none;
  background-color: #f0f0f0;
  cursor: pointer;
}
.filter-buttons button.active {
  background-color: #4caf50;
  color: white;
}
.conversations-list {
  margin-top: 20px;
}
.modal {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
}
.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 5px;
  width: 300px;
  text-align: center;
}
.close {
  position: absolute;
  top: 10px;
  right: 10px;
  cursor: pointer;
}
</style>