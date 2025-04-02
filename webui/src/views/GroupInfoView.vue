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
              <img v-if="member.photo && member.photo !== ''" :src="`data:image/jpeg;base64,${member.photo}`" alt="User Avatar" class="member-avatar" />                {{ member.username }}
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
/* Global Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Arial', sans-serif;
  background-color: #f5f6f7;
}

/* Header */
header {
  background-color: #0088cc;
  color: white;
  padding: 15px 20px;
  text-align: center;
  border-radius: 20px 20px 0 0;
}

header h1 {
  font-size: 24px;
  font-weight: bold;
}

/* Buttons */
button {
  font-size: 14px;
  padding: 10px 15px;
  border-radius: 15px;
  border: none;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  background-color: #0088cc;
  color: white;
  transition: background-color 0.3s;
}

button svg {
  margin-right: 8px;
  fill: white;
}

button:hover {
  background-color: #005f8a;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

/* Go Back Button */
.go-back-btn {
  margin: 20px 0;
  font-size: 16px;
  background-color: #4c6ef5;
}

.go-back-btn:hover {
  background-color: #3b5ef5;
}

/* Info View */
.info-view {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.username {
  flex-grow: 1;
}

.username-button {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.username-button button {
  background: none;
  color: #0088cc;
  padding: 5px;
}

.username-button span {
  font-size: 18px;
  font-weight: bold;
}

/* Profile Picture */
.pfp-view {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.edit-btn {
  background: none;
  border: none;
  color: #0088cc;
  font-size: 18px;
  cursor: pointer;
}

.profile-picture {
  flex-grow: 1;
  text-align: center;
}

.pfp-img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  margin-top: 10px;
}

/* Modal */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 20px;
  text-align: center;
  width: 300px;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}

.close {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 20px;
  cursor: pointer;
  color: #0088cc;
}

.modal h2 {
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: bold;
}

input[type="text"], input[type="file"] {
  padding: 10px;
  margin-top: 10px;
  width: 100%;
  border-radius: 15px;
  border: 1px solid #cccccc;
  outline: none;
  transition: border 0.3s;
}

input[type="text"]:focus, input[type="file"]:focus {
  border: 1px solid #0088cc;
}

button[type="button"] {
  margin-top: 10px;
}

/* Error Message */
.error-msg {
  color: #e74c3c;
  font-size: 14px;
  margin-top: 10px;
}

/* Members Section */
.members-section {
  padding: 15px;
  background-color: white;
  border-radius: 20px;
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}

.members-section h2 {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border-radius: 10px;
  background-color: #f1f1f1;
  cursor: pointer;
  transition: background-color 0.3s;
}

li:hover {
  background-color: #0088cc;
  color: white;
}

.member-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

/* Add and Leave Buttons */
.add-btn, .leave-btn {
  margin: 5px;
  font-size: 14px;
  background-color: #4c6ef5;
}

.add-btn:hover {
  background-color: #3b5ef5;
}

.leave-btn {
  background-color: #f14c4c;
}

.leave-btn:hover {
  background-color: #e03d3d;
}

</style>
