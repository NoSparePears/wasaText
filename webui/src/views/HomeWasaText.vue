<template>
  <div class="home-view">
    <header>
      
      <div class="header-buttons">
        <button @click="toggleSearchModal" class="button">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Chat
        </button>
        <button @click="getConversations" class="button">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"></use></svg>
          Refresh
        </button>
        <button @click="toggleGroupModal" class="button">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"></use></svg>
          New Group
        </button>
      </div>
    </header>
    
    <!-- Modal for searching users to start a new conversation -->
    <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="createNewChat" title="search">
      <template v-slot:header>
        <h3>Users</h3>
      </template>
    </Search>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <!-- Both Chat List -->
    <div class="chat-container">

      <div class="chat-section">
        <h3>Private Chats</h3>
        <ul v-if="conversations.length>0" class="conversation-list">
          <li v-for="conversation in conversations" :key="conversation.DestUserID" @click="openChat(conversation)" class="conversation-item">
    
            <img :src="conversation.destUser.photo" alt="User avatar" class="user-photo"/>
            <div class="chat-preview">
              <span class="username">{{ conversation.destUser.username }}</span>
              <span class="last-message-timestamp">{{ formatTimestamp(conversation.lastMessage.timestamp) }}</span>
              <span class="last-message">{{ conversation.lastMessage.content }}</span>
            </div>

          </li>
        </ul>
        <div v-else>No private chats. Start a new one!</div>
      </div>

      <div class="group-section">
        <h3>Groups</h3>
        <ul v-if="groups.length > 0" class="group-list">
          <li v-for="group in groups" :key="group.group.groupID" @click="openGroup(group.group)" class="group-item">
            <img :src="group.group.photo" alt="Group avatar" class="group-photo"/>
            <div class="group-preview">
              <span class="group-name">{{ group.group.groupName }}</span>
            </div>
          </li>
        </ul>
        <div v-else>No groups. Create a new one!</div>
      </div>
    

      
    </div>

    <!-- Modal for creating a new group -->
    <div v-if="groupModalVisible" class="modal">
      <div class="modal-content">
        <h3>Create New Group</h3>
      
        <!-- Group Name Input -->
        <input v-model="groupName" type="text" placeholder="Enter group name..." class="input-field"/>
      
        <!-- User Selection -->
        <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="addUserToGroup" title="search">
          <template v-slot:header>
            <h3>Select Users</h3>
          </template>
        </Search>

        <ul>
          <li v-for="user in selectedUsers" :key="user.id">
            {{ user.username }}
          </li>
        </ul>
      
        <button @click="toggleSearchModal" class="button">Add Users</button>
        <button @click="createGroup" class="button">Done</button>
        <button @click="toggleGroupModal" class="button close-button">Cancel</button>
      </div>
    </div>

    
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
      groupName: '',
      groupImage: null,
      selectedUsers: [],
      groups: [],
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
      
    },
    async getGroups() {
      this.errormsg = '';
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        let response = await this.$axios.get(`/profiles/${userID}/groups`, { 
          headers: { 'Authorization': token } });
        this.groups = response.data;
        if (!this.groups) this.groups = [];
      } catch (error) {
        this.errormsg = error.response?.data?.message || 'Error fetching conversations';
        console.error('Error fetching conversations:', error);
      }
    },
    async createGroup() {
      if (!this.groupName || this.selectedUsers.length === 0) {
        this.errormsg = 'Group name and at least one user are required.';
        return;
      }
      this.errormsg = '';
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        let response = await this.$axios.post(`/profiles/${userID}/groups`, {name: this.groupName, members: this.selectedUsers}, { 
          headers: { 'Authorization': token } 
          });
        
        if (!this.groups) this.groups = [];
        this.groups.push(response.data);  
        
        this.toggleGroupModal();  
        this.openGroup(response.data);
        
        
      } catch (error) {
        console.error("Error creating group:", error);
        this.errormsg = error.response?.data?.message || "Failed to create group.";
      }
      
      
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
    },
    async openGroup(group) {
      console.log('Opening group:');
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        this.$router.push({
          path: `/home/groups/${group.groupID}`,
          query: {
            name: group.groupName,
            groupID: group.groupID,
            avatar: group.photo,
          }
        });
      } catch (error) {
        this.errormsg = error.response?.data?.message || 'Error opening group';
        console.error('Error opening group:', error);
      }
    },
    addUserToGroup(user) {
      if (!this.selectedUsers.find(u => u.id === user.id)) {
        this.selectedUsers.push(user);
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
  },
  mounted() {
  // If user is not logged in, redirect to login page
  if (!sessionStorage.getItem('token')) {
    this.$router.push("/");
    return;
  }

  this.getConversations();
  this.getGroups();

  // Set up periodic fetching for conversations
  this.conversationInterval = setInterval(() => {
    this.getConversations();
  }, 1000);

  // Set up periodic fetching for groups
  this.groupInterval = setInterval(() => {
    this.getGroups();
  }, 1000);
},

beforeUnmount() {
  // Clear intervals when component is destroyed
  clearInterval(this.conversationInterval);
  clearInterval(this.groupInterval);
}

};
</script>

<style scoped>
.header-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
  padding: 10px;
}

.button {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #0088cc; /* Telegram blue */
  color: white;
  font-size: 14px;
  font-weight: bold;
  border: none;
  border-radius: 50px; /* Fully rounded */
  padding: 10px 15px;
  cursor: pointer;
  transition: background 0.3s ease, transform 0.1s ease-in-out;
  box-shadow: 0px 4px 6px rgba(0, 136, 204, 0.3);
}

.button:hover {
  background: #007bb5; /* Slightly darker on hover */
}

.button:active {
  transform: scale(0.95);
}


.conversation-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.chat-container {
  display: flex;
  justify-content: space-between;
  padding: 20px;
}
.chat-section, .group-section {
  width: 48%;
}

.conversation-item, .group-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 10px;
  background: #f1f1f1;
  margin: 5px 0;
  cursor: pointer;
  transition: background 0.3s;
}

.conversation-item:hover, .group-item:hover {
  background: #0088cc;
  color: white;
}

.user-photo, .group-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}
.chat-preview {
    display: flex;
    align-items: center;
    padding: 10px 0;
    cursor: pointer;
    transition: background 0.3s ease-in-out;
}

.username {
    font-weight: bold;
    color: #333;
    font-size: 16px;
    margin-right: 10px;
}

.last-message-timestamp {
    font-size: 14px;
    color: #555;
    margin-right: 10px;
}

.last-message {
    font-size: 14px;
    color: #555;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 200px; /* Optional: Adjust based on available space */
}
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
  border-radius: 10px;
  width: 400px;
  text-align: center;
}
.input-field, .file-input {
  display: block;
  width: 100%;
  margin: 10px 0;
  padding: 8px;
}
.close-button {
  background: red;
}

</style>
