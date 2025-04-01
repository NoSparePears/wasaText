<template>
    <header>
      <h1>Group Info</h1>
    </header>
    <!-- Button to go back to group -->
    <button @click="goBackToGroup" class="btn go-back-btn">
      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left"></use></svg>
      Go Back to Group
    </button>
    <div class="info-view">
        <!-- Sezione Name a sinistra -->
        <section class="username">
              <div class="username-button">
                    <button @click="openUsernameModal">
                      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-3"></use></svg>
                    </button>
                    <span class="text-lg font-semibold">{{ name }}</span>
              </div>
        </section>
    </div>
    
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
        
            <!-- Save Button -->
            <button @click="updatePhoto">Salva</button>
          </div>
        </div>
    </div>

    <section class="members-section">
            
        <button class="btn add-btn" @click="toggleSearchModal">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus"></use></svg>
            Add Member
        </button>
        <button class="btn leave-btn" @click="leaveGroup">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"></use></svg>
            Leave Group
        </button>

        <h2>Members</h2>
        <ul>
            <li v-for="member in members" :key="member.id">
                <img :src="member.avatar || 'default_avatar.jpg'" alt="User Avatar" class="member-avatar" />
                {{ member.username }}
            </li>
        </ul>

    </section>

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
            <input type="file" @change="handleFileChange">
            <button @click="updatePhoto">Salva</button>
      </div>
    </div>
    
    <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="addMember" title="search">
      <template v-slot:header>
        <h3>Select User</h3>
      </template>
    </Search>


</template>



<script>
import Search from '@/components/Search.vue';

export default {
    components: {
        Search
    },
    data() {
        return {
            errorMsg: "", // Messaggio di errore

            groupID: this.$route.query.groupID || '',
            name: this.$route.query.name || '',

            newUsername: '',
            isUsernameModalOpen: false,

            profilePicture: '',
            selectedFile: null, // Stores selected file
            isPhotoModalOpen: false, // Stato per il modale della foto

            members: [],
            searchModalVisible: false,
            
        }
        
    },
    methods: {
        async getMembers() {
            this.errormsg = '';
            const token = sessionStorage.getItem('token');
            const userID = sessionStorage.getItem('id');
            try {
                let response = await this.$axios.get(`/profiles/${userID}/groups/${this.groupID}/members`,
                {headers: { 'Authorization': token }
                })
                this.members = response.data;
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error fetching group members';
                console.error('Error fetching group members:', error);
            }
        },
        goBackToGroup() {
            console.log('Going back to group');
            try {
                this.$router.push({
                    path: `/home/groups/${this.groupID}`,
                    query: {
                        name: this.name || "Unknown",
                        groupID: this.groupID,
                        avatar: this.avatar || "default_propic.jpg",
                    }
                })
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error opening group s info';
                console.error('Error opening group info:', error);
            }
        },
        async getProfilePicture() {
          const userID = sessionStorage.getItem('id');
          const token = sessionStorage.getItem('token');
          try {
            // Fai una chiamata API per ottenere la foto del profilo
            const response = await this.$axios.get(`/profiles/${userID}/groups/${this.groupID}/g_photo`, 
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
        openPhotoModal() {
          this.isPhotoModalOpen = true;
        },
        // Chiudi il modale della foto
        closePhotoModal() {
          this.isPhotoModalOpen = false;
          this.selectedFile = null;
        },
        handleFileChange(event) {
          this.selectedFile = event.target.files[0]; // Store selected file
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
            const response = await this.$axios.put(`/profiles/${userID}/groups/${this.groupID}/g_photo`, formData, {
              headers: {
                "Authorization": token,
                "Content-Type": "multipart/form-data"
              }
            });

            this.getProfilePicture(); // Refresh profile picture

            this.closePhotoModal(); // Close modal after success
          } catch (error) {
            console.error("Errore nel caricamento della foto:", error);
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
        async updateGroupName() {
            this.errormsg = '';
            const token = sessionStorage.getItem('token');
            const userID = sessionStorage.getItem('id');
            try {
                await this.$axios.put(`/profiles/${userID}/groups/${this.groupID}/g_name`, { name: this.newUsername },
                {headers: { 'Authorization': token }
                });
                
                this.name = this.newUsername;

                this.closeUsernameModal();
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error updating group name';
                console.error('Error updating group name:', error);
            }
        },
        toggleSearchModal() {
            this.searchModalVisible = !this.searchModalVisible;
        },
        async addMember(user) {
            this.errormsg = '';
            const userID = sessionStorage.getItem('id');
            const token = sessionStorage.getItem('token');
            try {
                await this.$axios.put(`profiles/${userID}/groups/${this.groupID}/members`, {memberID: user.id},
                {headers: { 'Authorization': token }
                });
                this.getMembers();

                this.toggleSearchModal();
            } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error adding member to group';
                console.error('Error adding member to group:', error);
            }
        },
        async leaveGroup(){
          this.errormsg = '';
          const token = sessionStorage.getItem('token');
          const userID = sessionStorage.getItem('id');
          try {
            await this.$axios.delete(`/profiles/${userID}/groups/${this.groupID}`, 
                {headers: { 'Authorization': token }
                });
            this.$router.push({
                path: `/home`,
            })
          } catch (error) {
                this.errormsg = error.response?.data?.message || 'Error leaving group';
                console.error('Error leaving group:', error);
          }
        }
    },
    mounted() {
        if (!sessionStorage.getItem('token')) {
            this.$router.push("/");
            return;
        }
        this.getProfilePicture();
        this.getMembers();
        this.intervalId = setInterval(async () => {
          clearInterval(this.intervalId);
          await this.getMembers();
          this.intervalId = setInterval(this.getMembers, 5000);
        }, 1000);
    },
    beforeUnmount() {
        if (this.intervalId) {
          clearInterval(this.intervalId);
        }
    }
}

</script>

<style scoped>
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
.close {
  position: absolute;
  top: 10px;
  right: 15px;
  font-size: 24px;
  cursor: pointer;
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
  animation: fadeIn 0.3s ease-in-out;
}

.modal-content {
  background: #ffffff;
  border-radius: 15px;
  padding: 20px;
  width: 320px;
  box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.2);
  text-align: center;
  animation: slideIn 0.3s ease-in-out;
}

.modal-content h3 {
  margin-bottom: 15px;
  font-size: 18px;
  color: #0088cc;
}

.modal input {
  width: 100%;
  padding: 10px;
  border-radius: 10px;
  border: 1px solid #cccccc;
  outline: none;
  transition: border 0.3s;
  font-size: 16px;
}

.modal input:focus {
  border: 1px solid #0088cc;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: 15px;
}

.btn {
  padding: 10px 15px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 16px;
  transition: background 0.3s, transform 0.2s;
  width: 48%;
}

.save-btn {
  background: #0088cc;
  color: white;
  border: none;
}

.save-btn:hover {
  background: #0077b6;
  transform: scale(1.05);
}

.cancel-btn {
  background: #f1f1f1;
  color: #333;
  border: none;
}

.cancel-btn:hover {
  background: #e0e0e0;
  transform: scale(1.05);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(-10px); }
  to { transform: translateY(0); }
}
</style>
