<template>
  <header>
    <h1>Profilo</h1>
  </header>
  <div class="profile-view">
    <!-- Sezione Username a sinistra -->
    <section class="username">
      <div class="username-button">
        <button @click="openUsernameModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-3"></use></svg>
        </button>
        <span class="text-lg font-semibold">{{ username }}</span>
      </div>
    </section>

    <div class="pfp-view">
      <!-- Sezione Immagine profilo a destra -->
      <section class="profile-picture">
        <img :src="`data:image/jpg;base64,${profilePicture}`" class="profile-pic">        
      </section>
      <div class="pfp-button">
        <button @click="openPhotoModal">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2"></use></svg>
        </button>
      </div>
    </div>
  </div>

  <!-- Modale per modificare username -->
  <div v-if="isUsernameModalOpen" class="modal">
    <div class="modal-content">
      <span @click="closeUsernameModal" class="close">&times;</span>
      <h2>Modifica Username</h2>
      <input v-model="newUsername" type="text" placeholder="Nuovo Username">
      <button @click="updateUsername">Salva</button>
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
    </div>
  </div>

  <!-- Modale per modificare foto del profilo -->
  <div v-if="isPhotoModalOpen" class="modal">
    <div class="modal-content">
      <span @click="closePhotoModal" class="close">&times;</span>
      <h2>Modifica Foto Profilo</h2>
      <input type="file" @change="previewPhoto">
      <button @click="updatePhoto">Salva</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      profilePicture: sessionStorage.photo, // Percorso di default
      username: sessionStorage.username, // Username dell'utente dal database
      newUsername: "", // Nuovo nome utente da inviare
      isUsernameModalOpen: false, // Stato per il modale dello username
      isPhotoModalOpen: false, // Stato per il modale della foto
      errorMsg: "" // Messaggio di errore
    };
  },
  methods: {
    // Apri il modale per modificare lo username
    openUsernameModal() {
      this.isUsernameModalOpen = true;
      this.newUsername = this.username; // Pre-compila il campo con lo username attuale
    },
    // Chiudi il modale dello username
    closeUsernameModal() {
      this.isUsernameModalOpen = false;
      window.location.reload(); //reloado la pagina
    },
    // Funzione per aggiornare lo username
    async updateUsername() {
      if (this.newUsername !== this.username) {
        try {
          // Fai una chiamata API per aggiornare lo username nel database
          
          let _ = await this.$axios.put(`/profiles/${sessionStorage.id}/username`, { username: this.newUsername }, { headers: { 'Authorization': `${sessionStorage.token}` } })
          // Assegna il nuovo username alla variabile username per l'aggiornamento della pagina;

          // Aggiorna lo username nel frontend
          sessionStorage.username = this.newUsername; // Salva nel sessionStorage
          this.closeUsernameModal(); // Chiudi il modale
        } catch (error) {
          this.errorMsg = "Errore di rete!";
        }
      }
    },
    
    // Apri il modale per modificare la foto del profilo
    openPhotoModal() {
      this.isPhotoModalOpen = true;
    },
    // Chiudi il modale della foto
    closePhotoModal() {
      this.isPhotoModalOpen = false;
    },
    // Funzione per caricare l'anteprima della foto
    previewPhoto(event) {
      const file = event.target.files[0];
      if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
          this.profilePicture = e.target.result; // Mostra l'anteprima della nuova foto
        };
        reader.readAsDataURL(file);
      }
    },
    // Funzione per aggiornare la foto del profilo
    async updatePhoto() {
      try {
        const formData = new FormData();
        formData.append('photo', this.$refs.fileInput.files[0]);

        const response = await fetch(`/profiles/${sessionStorage.userID}/photo`, {
          method: 'PUT',
          body: formData,
        });

        if (response.ok) {
          const data = await response.json();
          this.profilePicture = data.newPhotoUrl; // Aggiorna l'URL della foto del profilo
          sessionStorage.setItem('photo', this.profilePicture); // Salva nel sessionStorage
          this.closePhotoModal(); // Chiudi il modale
        } else {
          const error = await response.json();
          this.errorMsg = error.message || "Errore durante l'aggiornamento della foto!";
        }
      } catch (error) {
        this.errorMsg = "Errore di rete!";
      }
    }
  }
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

/* Stili per il modale */
.modal {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 5px;
  width: 300px;
  text-align: center;
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  cursor: pointer;
}

/* Messaggio di errore */
.error-msg {
  color: red;
  font-size: 14px;
}
</style>

