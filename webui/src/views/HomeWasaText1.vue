<template>
  <div class="home-view">
    <header>
      <h1>Chat</h1>
      <div class="header-buttons">
        <button @click="searchModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Chat
        </button>
        <button @click="groupModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Group
        </button>
        <button @click="getConversations">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"></use></svg>
          Refresh
        </button>
      </div>
    </header>
    
    <!-- Modale utilizzato per la creazione di un nuovo gruppo -->
    <Group :show="groupModalVisible" @close="groupModal" title="search">
      <template v-slot:header>
        <h3>Select users</h3>
      </template>
    </Group>
    <!-- Modale utilzzato per la ricerca degli utenti con cui aprire una nuova conversazione -->
    <Search :show="searchModalVisible" @close="searchModal" title="search">
      <template v-slot:header>
        <h3>Users</h3>
      </template>
    </Search>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div v-if="conversations.length === 0">
      Create new chat
    </div>
    <ul v-else>
      <li v-for="conversation in conversations" :key="conversation.id">
        {{ conversation.name }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  data() {
    return {
      searchModalVisible: false,
      groupModalVisible: false,
      conversations: [],
      errormsg: ''
    };
  },
  methods: {
    searchModal() {
      this.searchModalVisible = !this.searchModalVisible;
    },
    groupModal() {
      this.groupModalVisible = !this.groupModalVisible;
    },
    async getConversations() {
      try {
        //effettuo richiesta api
        let response = await this.$axios.get(`/profiles/${sessionStorage.id}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
        //salvo dati conversazioni
        this.conversations = response.data
      } catch (error) {
        console.error('Errore nel recupero delle conversazioni:', error);
      }
    }
  },
  async createNewChat(user) {
      try {
        let response = await this.$axios.post(`/profiles/${sessionStorage.id}/conversations`, { userId: user.id }, { headers: { 'Authorization': sessionStorage.token } });
        this.conversations.push(response.data);
        this.toggleSearchModal();
      } catch (error) {
        this.errormsg = 'Error creating new chat';
        console.error('Error creating new chat:', error);
      }
    },

  created() {
    this.getConversations();
  }
};
</script>

<style>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}
.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
}
.close {
  cursor: pointer;
  font-size: 24px;
  position: absolute;
  right: 10px;
  top: 10px;
}
</style>