<!-- 

Modale utilizzato per la ricerca di utenti

-->

<script>

export default {
  // Props passati al componente
  props: {
    show: Boolean,  // Utilizzato per mostrare o nascondere il modal
    users: Array, // Lista di utenti da filtrare
    title: String,  // Titolo del modale
  },
  data() {
    return {
      errorMsg: "",

      // Variabile utilizzate per la ricerca degli utenti
      searchText: "",

      // Utilizzato per la verifica dell'username inserito
      usernameValidation: new RegExp('^\\w{0,16}$'),

      // Username dell'utente che ha effettuato il login
      username: sessionStorage.username,

      // Lisat di utenti filtrati in base all'username inserito
      filteredUsers: [],
    };
  },
  methods: {
    // Funzione utilizzata per chiudere il modale salvando in localStorage l'utente selezionato con cui aprire un conversazione
    closeModal() {
      this.searchText = "";
      window.location.reload();
      this.$emit('close');
    },
    // Funzione utilizzata per filtrare gli utenti in base all'username inserito
    async filterUsers() {
      this.errorMsg = "";
      // Lista degli utenti
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
    // Funzione utilizzara per selezionare l'utente con cui aprire una conversazione
    async selectUser(userToSend, username, photo) {
      localStorage.clear();
      localStorage.userID = userToSend;
      localStorage.username = username;
      localStorage.photo = photo;
      this.$router.push('/conversation');
      closeModal();
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
              <!-- Input per l'inserimento dell'utente da cercare -->
              <div class="search-input">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <!-- Risultato della ricerca -->
              <div class="search-results">
                <div v-for="user in filteredUsers" :key="user.userId"
                  @click="selectUser(user.userId, user.username, user.photo)">
                    <div class="user">
                      <p v-if="user.username != username" >{{ user.username }}</p>
                    </div>
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