<!-- In questa pagina l'utente effettua il login -->

<template>
    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
    <div class="login-container">
      <form @submit.prevent="doLogin">
        <h1>WasaText</h1>
        <input type="text" v-model="username" placeholder="Enter your username" />
        <button type="submit">Login</button>
      </form>
    </div>
  </template>

<script>
    export default {
  data() {
    return {
      // Username input dell'utente che si sta loggando
      username: "",
      errorMsg: "",

      // Verifica per il campo username
      usernameValidation: new RegExp('^\\w{3,16}$'),
    }
  },
  emits: ['login-success'],
  methods: {
    // funzione per effettuare login
    async doLogin() {
        try {
            // faccio rischiesta login
            let response = await this.$axios.post('/session', {
                username: this.username
            });

            sessionStorage.id = response.data.id;
            sessionStorage.username = this.username;
            sessionStorage.token = response.data.id;
            sessionStorage.photo = response.data.photo;
            // Reindirizza l'utente alla home
            this.$router.push('/home');
            // Emette l'evento di login avvenuto con successo
            this.$emit('login-success');
        } catch (e) {
            this.errorMsg = e.toString();
        }
    },
    mounted() {
        // Se l'utente è già loggato, reindirizza alla home
        if (sessionStorage.token) {
          this.$router.push("/home");
          return;
        }
        // Altrimewnti cancella i dati dell'utente dalla sessionStorage
        sessionStorage.clear();
    },
    }
}
</script>