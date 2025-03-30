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

      <!-- Edit Button on the Left -->
      <button @click="openPhotoModal" class="edit-btn">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2"></use></svg>
      </button>
      <!-- Profile Picture Section ["/feather-sprite-v4.29.0.svg#edit-2"]-->
      <section class="profile-picture">
        <img v-if="profilePicture" :src="profilePicture" alt="Profile Picture" class="pfp-img" />
      </section>
  
      <!-- Modal for Updating Profile Picture -->
      <div v-if="isPhotoModalOpen" class="modal">
        <div class="modal-content">
          <span @click="closePhotoModal" class="close">&times;</span>
          <h2>Modifica Foto Profilo</h2>

          <!-- Image Preview -->
          <div v-if="previewImage">
            <img :src="previewImage" alt="Anteprima" class="preview-img" />
          </div>

          <!-- File Upload Input -->
          <input type="file" @change="previewPhoto" accept="image/*">

          <!-- Save Button -->
          <button @click="updatePhoto">Salva</button>
        </div>
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
      errorMsg: "", // Messaggio di errore

      username: sessionStorage.username, // Username dell'utente dal database
      newUsername: "", // Nuovo nome utente da inviare
      isUsernameModalOpen: false, // Stato per il modale dello username
      
      profilePicture: '', // Foto del profilo
      previewImage: null, // Stores preview image before upload
      selectedFile: null, // Stores selected file
      isPhotoModalOpen: false, // Stato per il modale della foto
      
    };
  },
  methods: {
    async getProfilePicture() {
      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');
      try {
        // Fai una chiamata API per ottenere la foto del profilo
        const response = await this.$axios.get(`/profiles/${userID}/photo`, 
          { headers: { 'Authorization': token } 
        });
        // Extract base64 data from response
        if (response.data && response.data.profile_picture) {
          this.profilePicture = `data:image/jpeg;base64,${response.data.profile_picture}`;
        } else {
          console.error("Profile picture data is missing in the response");
        }
      } catch (error) {
        console.error("Errore nel recupero della foto del profilo:", error);
      }
    },
    // Apri il modale per modificare lo username
    openUsernameModal() {
      this.isUsernameModalOpen = true;
      this.newUsername = this.username; // Pre-compila il campo con lo username attuale
    },
    // Chiudi il modale dello username
    closeUsernameModal() {
      this.isUsernameModalOpen = false;
    },
    // Funzione per aggiornare lo username
    async updateUsername() {
      if (this.newUsername !== this.username) {
        try {
          // Fai una chiamata API per aggiornare lo username nel database
          await this.$axios.put(`/profiles/${sessionStorage.getItem('id')}/username`, { username: this.newUsername }, { headers: { 'Authorization': sessionStorage.getItem('token') } })
          // Aggiorna lo username nel frontend
          sessionStorage.username = this.newUsername; // Salva nel sessionStorage
          this.username = this.newUsername
          this.closeUsernameModal(); // Chiudi il modale
        } catch (error) {
          if (error.response && error.response.data) {
            this.errorMsg = error.response.data; // Set error message from the backend
          } else {
            this.errorMsg = "Errore di rete!";
          }
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
      this.previewImage = null;
      this.selectedFile = null;
    },
    // Preview selected image
    previewPhoto(event) {
      const file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        const reader = new FileReader();
        reader.onload = (e) => {
          this.previewImage = e.target.result; // Show preview
        };
        reader.readAsDataURL(file);
      }
    },
    // Upload and update profile picture
    async updatePhoto() {
      if (!this.selectedFile) {
        alert("Seleziona un'immagine prima di salvare!");
        return;
      }

      const userID = sessionStorage.getItem('id');
      const token = sessionStorage.getItem('token');

      const formData = new FormData();
      formData.append("profile_picture", this.selectedFile);

      try {
        const response = await this.$axios.put(`/profiles/${userID}/photo`, formData, {
          headers: {
            "Authorization": token,
            "Content-Type": "multipart/form-data"
          }
        });

        this.getProfilePicture();
        
        this.closePhotoModal(); // Close modal after success
      } catch (error) {
        console.error("Errore nel caricamento della foto:", error);
      }
    }
  },
  mounted() {
      // Se l'utente non è loggato, reindirizza alla pagina di login
      if (!sessionStorage.getItem('token')) {
        this.$router.push("/");
        return;
      }
      this.getProfilePicture();
      
    },
};
</script>

<style>
body {
  font-family: Arial, sans-serif;
  background-color: #f3f4f6;
  color: #333;
}
.header {
  background: #ffffff;
  padding: 15px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  text-align: center;
}
.profile-view {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: white;
  margin: 20px auto;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}
.username-button {
  display: flex;
  align-items: center;
  gap: 10px;
}
.icon-button {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 5px;
}
.pfp-view {
  display: flex;
  align-items: center;
  gap: 10px; /* Space between profile pic & button */
}

.profile-picture {
  width: 100px;  /* Circle size */
  height: 100px;
  border-radius: 50%; /* Makes it a circle */
  overflow: hidden; /* Prevents image from spilling out */
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f0f0; /* Background color (optional) */
}

.pfp-img {
  width: 120%; /* Ensures full coverage */
  height: 120%;
  object-fit: cover; /* Ensures image fits without stretching */
}

/* Edit Button */
.edit-btn {
  display: block;
  margin-top: 10px;
  background-color: #007bff;
  color: white;
  padding: 5px 10px;
  border: none;
  cursor: pointer;
  border-radius: 5px;
}
.edit-btn:hover {
  background-color: #0056b3;
}

/* Modal Styling */
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
  text-align: center;
  max-width: 400px;
}

.close {
  position: absolute;
  top: 10px;
  right: 15px;
  font-size: 24px;
  cursor: pointer;
}

/* Preview Image */
.preview-img {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  margin-top: 10px;
}
.input-field, .input-file {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ddd;
  border-radius: 5px;
}
.primary-button {
  background: #0088cc;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
.primary-button:hover {
  background: #0077b3;
}
</style>

