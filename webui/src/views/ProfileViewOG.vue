<template>
  <header>
          <h1>Profilo</h1>
  </header>
  <div class="profile-view">
          <!-- Sezione Username a sinistra -->
          <section class="username">
              <div class="username-button">
                  <button @click="setUsername">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-3"></use></svg>
                  </button>
                  <span class="text-lg font-semibold">{{ username }}</span>
              </div>
          </section>
          <div class = "pfp-view">
    <!-- Sezione Immagine profilo a destra -->
    <section class="profile-picture">
      <img :src="`data:image/jpg;base64,${profilePicture}`" class="profile-pic">        
    </section>
    <div class="pfp-button">
      <button @click="setMyPhoto">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2"></use></svg>
      </button>
    </div>
  </div>
  </div>
  
</template>

<script>
export default {
  data() {
    return {
      profilePicture: sessionStorage.photo, // Percorso di default
      username: sessionStorage.username, // Username dell'utente dal database
      errorMsg: ""
    };
  },
  methods: {
    
    //funzione per cambiare username "/profiles/:userID/username", rt.wrap(rt.setMyUserName, true)
    
    //funzione per cambiare foto "/profiles/:userID/photo", rt.wrap(rt.setMyPhoto, true)
  },
 
};
</script>

<style scoped>
/* Imposta la struttura a colonna per evitare problemi di posizionamento */
.profile-view {
 
  text-align: left;
  display: flex;
  flex-direction: column-reverse; /* Organizza gli elementi in colonna */
  align-items: left; /* Allinea tutto al centro */
}

.pfp-button {
  position: relative;
  left: 120px;
}
.pfp-button button {
  background: none;
  border: none;
  cursor: pointer; 
}

.profile-pic {
  width: 100px; /* Dimensione regolabile */
  height: 100px;
  border-radius: 50%; /* Trasforma in cerchio */
  object-fit: cover; /* Ritaglia l'immagine senza distorsioni */
  border: 2px solid #ccc; /* Bordo opzionale */
}

.username {
  text-align: left;
  margin-top: 50px;
}

.username-button {
  display: flex;
  align-items: center;
  gap: 5px; /* Spazio tra bottone e testo */
}

.username-button button {
  background: none;
  border: none;
  cursor: pointer;
}

</style>
