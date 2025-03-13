<template>
  <div v-if="show" class="modal">
    <div class="modal-content">
      <span class="close" @click="$emit('close')">&times;</span>
      <slot name="header"></slot>
      <input 
        type="text" 
        v-model="query" 
        @input="onInput" 
        placeholder="Search users..." 
        class="search-input"
      />
      <ul class="user-list">
        <li 
          v-for="destUser in users" 
          :key="destUser.id" 
          @click="selectUser(destUser)"
          class="user-item"
        >
          <span v-html="highlightMatch(destUser.username)"></span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import debounce from 'lodash.debounce';

export default {
  props: {
    show: Boolean,
    title: String
  },
  data() {
    return {
      query: '',
      users: []
    };
  },
  watch: {
    query(newQuery) {
      if (newQuery.length > 2) {
        this.debouncedSearch();
      } else {
        this.users = [];
      }
    }
  },
  created() {
    this.debouncedSearch = debounce(this.searchUsers, 300);
  },
  methods: {
    async searchUsers() {
      try {
        const url = `/profiles?username=${this.query}`;
        let response = await this.$axios.get(url, {
          headers: { 'Authorization': `${sessionStorage.getItem('token')}` }
        });
        this.users = response.data;
      } catch (error) {
        console.error('Error searching users:', error);
      }
    },
    selectUser(destUser) {
      this.$emit('user-selected', destUser);
    },
    highlightMatch(username) {
      if (!this.query) return username;
      const regex = new RegExp(`(${this.query})`, 'gi');
      return username.replace(regex, '<span class="highlight">$1</span>');
    }
  }
};
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: #ffffff;
  border-radius: 20px;
  padding: 20px;
  width: 300px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
  animation: fadeIn 0.3s ease-in-out;
}

.close {
  float: right;
  font-size: 24px;
  cursor: pointer;
  color: #0088cc;
}

.search-input {
  width: 100%;
  padding: 10px;
  border-radius: 15px;
  border: 1px solid #cccccc;
  outline: none;
  transition: border 0.3s;
}

.search-input:focus {
  border: 1px solid #0088cc;
}

.user-list {
  list-style: none;
  padding: 0;
  margin: 10px 0 0;
}

.user-item {
  padding: 10px;
  background: #f1f1f1;
  border-radius: 10px;
  margin: 5px 0;
  cursor: pointer;
  transition: background 0.3s;
}

.user-item:hover {
  background: #0088cc;
  color: white;
}

.highlight {
  background-color: yellow;
  font-weight: bold;
}

@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.9); }
  to { opacity: 1; transform: scale(1); }
}
</style>
