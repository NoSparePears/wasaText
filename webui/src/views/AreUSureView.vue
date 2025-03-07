<template>
    <header>
      <h1>Sei sicuro?</h1>
    </header>
    <div class="buttons">
      <button @click="logOut">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#check-circle"></use></svg>
      </button>
      <button @click="noLogOut">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#slash"></use></svg>
      </button>
    </div>
  </template>

<script>
export default {
  methods: {
    // Funzione utilizzata per il logout dell'utente
    logOut() {
      sessionStorage.clear();
      sessionStorage.isLoggedIn = false;
      this.$router.push("/");
    },
    noLogOut() {
      // Se l'utente non vuole fare il log out, rimandalo alla pagina principale
      this.$router.push('/home');
    }
  }
};
</script>

<style scoped>
header {
  text-align: center; /* Centra la scritta */
  padding: 20px;
}

h1 {
  margin: 0;  /* Rimuove il margine di default */
  font-size: 40px; /* Imposta la dimensione del testo */
}

.buttons {
  display: flex;
  justify-content: center; /* Centra i bottoni orizzontalmente */
  gap: 30px; /* Aggiunge spazio tra i bottoni */
  margin-top: 20px;
}

button {
  background: none;
  border: none;
  cursor: pointer;
}

button svg {
  width: 50px; /* Rende le icone molto più grandi */
  height: 50px;
}
</style>

  