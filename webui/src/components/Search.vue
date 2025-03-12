<template>
  <div v-if="show" class="modal">
    <div class="modal-content">
      <span class="close" @click="$emit('close')">&times;</span>
      <slot name="header"></slot>
      <input type="text" v-model="query" @input="onInput" placeholder="Search users..." />
      <ul>
        <li v-for="destUser in users" :key="destUser.id" @click="selectUser(destUser)">
          {{ destUser.username }}
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
    }
  }
};
</script>

<style>
.custom-link {
  color: inherit;
  /* This will make the link have the same color as the surrounding text */
  text-decoration: none;
  /* This will remove the underline */
}

.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: table;
  transition: opacity 0.3s ease;
}

.modal-wrapper {
  display: table-cell;
  vertical-align: middle;
}

.modal-container {
  width: 350px;
  margin: 0px auto;
  background-color: #fff;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-header {
  height: 70px;
  padding: 20px 15px 10px 15px;
}

.modal-header h3 {
  margin-top: 0;
  font-size: 25px;
  color: #42b983;
}

.modal-header button {
  color: rgb(86, 86, 86);
  background: none;
  border: none;
  padding: 5px;
  line-height: 12px;
  font-size: 15px;
}

.modal-header button svg {
  width: 20px;
  height: 20px;
}


.search-input {
  padding: 0 15px;
}

.search-input input {
  height: 30px;
  width: 100%;
  outline: none;
  border-radius: 3px;
  border: 1px solid rgb(179, 179, 179)
}

.search-results {
  font-size: 15px;
  padding: 10px 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  max-height: 200px;
  overflow-y: scroll;
}

.modal-default-button {
  float: right;
}

.username-form {
  display: flex;
  flex-direction: column;
  padding: 0 15px;
}

.username-form input {
  margin-bottom: 10px;
  margin-top: 5px;
  outline: none;
  border-radius: 3px;
  border: 1px solid rgb(179, 179, 179)
}

.username-form button {
  margin-bottom: 15px;
}
</style>