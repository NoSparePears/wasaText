<!-- 

Modale utilizzato per la creazione di un gruppo

L'utente loggato può:
- selezionare un nome per il gruppo
- cercare un utente specifico per aggiungerlo al gruppo
- visualizzare la lista di utenti aggiunti
- creare un gruppo con gli utenti selezionati

-->

<script>
export default {
  props: {
    show: Boolean,
    users: Array,
    title: String,
  },
  data() {
    return {
      errorMsg: "",

      // Utilizzato per la ricerca degli utentei da aggiungere al gruppo
      searchText: "",

      // Utilizzate per dare il nome al gruppo da creare
      groupName: "",

      // Utilizzato per la verifica del nome del gruppo e dell'username inserito
      usernameValidation: new RegExp('^\\w{0,16}$'),

      // Lista degli utenti cercati
      filteredUsers: [],

      // Username dell'utente loggato
      owner: sessionStorage.username,
      selectedUsers: [], // Lista di utenti selezionati
    };
  },
  methods: {
    // Funzione utilizzata per chiudere il modale
    closeModal() {
      this.searchText = "";
      this.groupName = "";
      this.selectedUsers = [];
      this.$emit('close');
    },
    // Funzione utilizzata per la ricerca di utenti da aggiungere al gruppo
    async filterUsers() {
      this.errorMsg = "";
      this.filteredUsers = this.users;

      if (this.searchText.length > 0) {
        if (this.searchText.length > 16 || !this.usernameValidation.test(this.searchText)) {
          this.errorMsg = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
          this.filteredUsers = [];
          return;
        }

        if (this.title === "search") {
          try {
            // Effettua una richiesta GET al server per ottenere gli utenti in base alla ricerca effettuata
            const url = `/profiles?username=${this.searchText}`;
            let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
            // In base al risultato della GET assegna la lista di utenti filtrati
            if (response.data == null) {
              this.filteredUsers = [];
              return;
            }
            this.filteredUsers = response.data;
          } catch (e) {
            this.errorMsg = e.toString();
            this.filteredUsers = [];
          }
        } else {
          this.filteredUsers = this.users.filter(user => user.username.toLowerCase().includes(this.searchText.toLowerCase()));
        }
      }
    },
    // Funzione utilizzata per creare un gruppo con gli utenti selezionati
    async createGroup() {
      // Controlla se il nome del gruppo è valido
      if (this.groupName.length < 3 || this.groupName.length > 16) {
        this.errorMsg = "Invalid group name, it must contains min 3 characters and max 16 characters";
        return;
      }
      try {
        // Effettua una richiesta POST al server per creare un gruppo passando il nome del gruppo e la lista degli utenti
        let response = await this.$axios.post(`/profiles/${sessionStorage.userID}/groups`, {
          groupName: this.groupName,
          users: this.selectedUsers,
        }, { headers: { 'Authorization': `${sessionStorage.token}` } });
        // Assegna al localStorage l'id del gruppo e il nome
        localStorage.clear();
        localStorage.userID = response.data.group.GroupId;
        localStorage.username = response.data.group.groupName;
        localStorage.photo = response.data.group.photo;
        // Chiude il modale
        this.closeModal();
        // Reinderizza alla pagina della conversazione
        this.$router.push(`/conversations/${response.data.conversation.conversationId}`);
      } catch (e) {
        this.errorMsg = e.toString();
      }
    },
    // Seleziona l'utente da aggiungere alla lista degli utenti del gruppo
    selectUser(user) {
      // Controlla se l'utente è già nella lista
      if (!this.selectedUsers.find(u => u.username === user.username)) {
        this.selectedUsers.push(user); // Aggiungi utente selezionato
      }
    },
    // Rimuove un utente dalla lista degli utenti del gruppo
    removeUser(username) {
      // Rimuove l'utente dalla lista selezionata
      this.selectedUsers = this.selectedUsers.filter(user => user.username !== username);
    },
  },
  watch: {
    searchText() {
      this.filterUsers();
    },
    show() {
      this.filteredUsers = this.users;
    }
  },
}
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <slot name="header">default header</slot>
            <button class="like-btn" @click="closeModal">
              <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#x" />
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <slot name="body">
              <!-- Selezione del nome del gruppo -->
              <div class="search-input">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="text" v-model="groupName" placeholder="Select group name" />
              </div>
              <!-- Campo di ricerca -->
              <div class="search-input">
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <p></p>
              <div class="btn-group me-2">
                <button class="btn btn-sm btn-outline-primary" @click="createGroup">Create Group</button>
              </div>

              <!-- Risultati della ricerca -->
              <div class="search-results">
                <div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
                  <p>{{ user.username }}</p>
                </div>
              </div>

              <!-- Lista di utenti selezionati -->
              <div class="selected-users">
                <h4>Selected Users:</h4>
                <span class="selected-user">{{ owner }}</span>
                <div v-for="user in selectedUsers" :key="user.userId" class="selected-user">
                  <span>{{ user.username }}</span>
                  <button v-if="user.username !== owner" @click="removeUser(user.username)">Remove</button>
                </div>
              </div>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>


<style>
.selected-users {
  margin-top: 20px;
  padding: 10px;
  border-top: 1px solid #ccc;
}

.selected-user {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.selected-user span {
  font-size: 14px;
  font-weight: bold;
}

.selected-user button {
  background: red;
  color: white;
  border: none;
  border-radius: 5px;
  padding: 5px 10px;
  cursor: pointer;
}
</style>