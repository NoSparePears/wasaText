<template>
  <div class="home-view">
    <header>
      <h1>Chat</h1>
      <div class="header-buttons">
        <button @click="toggleSearchModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Chat
        </button>
        <button @click="toggleGroupModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Group
        </button>
        <button @click="getConversations">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"></use></svg>
          Refresh
        </button>
      </div>
    </header>
    
    <!-- Modal for creating a new group -->
    <Group :show="groupModalVisible" @close="toggleGroupModal" title="search">
      <template v-slot:header>
        <h3>Select users</h3>
      </template>
    </Group>
    <!-- Modal for searching users to start a new conversation -->
    <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="createNewChat" title="search">
      <template v-slot:header>
        <h3>Users</h3>
      </template>
    </Search>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div v-if="conversations.length === 0">
      Create new chat
    </div>
    <ul v-else>
      <li v-for="conversation in conversations" :key="conversation.DestUserID" @click="openChat(conversation)">
        {{ conversation.destUser.username}}
      </li>
    </ul>
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
      searchModalVisible: false,
      groupModalVisible: false,
      conversations: [],
      errormsg: ''
    };
  },
  methods: {
    toggleSearchModal() {
      this.searchModalVisible = !this.searchModalVisible;
    },
    toggleGroupModal() {
      this.groupModalVisible = !this.groupModalVisible;
    },
    async getConversations() {
      this.errormsg = '';
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        let response = await this.$axios.get(`/profiles/${userID}/conversations`, { 
          headers: { 'Authorization': token } });
        this.conversations = response.data;
        if (!this.conversations) this.conversations = [];
      } catch (error) {
        this.errormsg = error.response?.data?.message || 'Error fetching conversations';
        console.error('Error fetching conversations:', error);
      }
    },
    async createNewChat(destUser) {
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        // Refresh conversation list first
        await this.getConversations();
    
        console.log('Current conversations:', this.conversations);
        console.log('Checking for existing conversation with:', destUser.id);    
        // Check if the conversation with destUser already exists
        const existingConvo = this.conversations.find(conversation => conversation.destUser.id === destUser.id);
        if (existingConvo) {
          console.log('Existing conversation found:', existingConvo);
          this.openChat(existingConvo);
          return;
        }
  
        console.log('No existing conversation, creating a new one...');
      
        let response = await this.$axios.get(`/profiles/${userID}/conversations/${destUser.id}`, { 
          headers: { 'Authorization': token } 
          });
  
        if (!this.conversations) this.conversations = [];
          this.conversations.push(response.data);
      
          this.openChat(response.data);
          this.toggleSearchModal();
  
        } catch (error) {
          this.errormsg = error.response?.data?.message || 'Error creating new chat';
          console.error('Error creating new chat:', error);
        }
        this.$router.push({
          path: `/home/${destUser.id}`,
          query: {
              name: destUser.username,
              destID: destUser.id,
              avatar: destUser.photo
          }
        })
    },
    async openChat(conversation) {
      console.log('Opening chat:');
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        this.$router.push({
          path: `/home/${conversation.destUser.id}`,
          query: {
            name: conversation.destUser.username,
            destID: conversation.destUser.id,
            avatar: conversation.destUser.photo
          }
        });
      } catch (error) {
        this.errormsg = error.response?.data?.message || 'Error opening chat';
        console.error('Error opening chat:', error);
      }
    }
  },
  mounted() {
    // Se l'utente non è loggato, reindirizza alla pagina di login
    if (!sessionStorage.getItem('token')) {
      this.$router.push("/");
      return;
    }
    this.getConversations();
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