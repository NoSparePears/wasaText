<template>
    <div class="group-info-view">
        <!-- Header -->
        <header class="group-header">
            <button class="back-button" @click="goBackToGroup">←</button>
            <div class="group-avatar-container">
                <img :src="avatar" alt="Group Avatar" class="group-avatar" />
                <button class="edit-btn" @click="toggleAvatarModal">✎</button>
            </div>
            <div class="group-name-container">
                <h1>{{ name }}</h1>
                <button class="edit-btn" @click="toggleNameModal">✎</button>
            </div>
        </header>

        <!-- Members List Section -->
        <section class="members-section">
            <h2>Members</h2>
            <ul>
                <li v-for="member in members" :key="member.id">
                    <img :src="member.avatar || 'default_avatar.jpg'" alt="User Avatar" class="member-avatar" />
                    {{ member.username }}
                </li>
            </ul>
            <button class="btn add-btn" @click="toggleSearchModal">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus"></use></svg>
                Add Member
            </button>
            <button class="btn leave-btn" @click="leaveGroup">Leave Group
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"></use></svg>

            </button>
        </section>

        <!-- Edit Name Modal -->
        <div v-if="showNameModal" class="modal">
            <div class="modal-content">
                <h3>Edit Group Name</h3>
                <input type="text" v-model="newGroupName" placeholder="Enter new name" />
                <div class="modal-buttons">
                    <button class="btn save-btn" @click="updateGroupName">Save</button>
                    <button class="btn cancel-btn" @click="toggleNameModal">Cancel</button>
                </div>
            </div>
        </div>

        <!-- Edit Avatar Modal -->
        <div v-if="showAvatarModal" class="modal">
            <div class="modal-content">
                <h3>Change Group Avatar</h3>
                <input type="file" @change="handleAvatarUpload" />
                <div class="modal-buttons">
                    <button class="btn save-btn" @click="updateGroupAvatar">Save</button>
                    <button class="btn cancel-btn" @click="toggleAvatarModal">Cancel</button>
                </div>
            </div>
        </div>
        <!-- Modal for searching users to add a new member -->
        <Search :show="searchModalVisible" @close="toggleSearchModal" @user-selected="addMember" title="search">
          <template v-slot:header>
            <h3>Select User</h3>
          </template>
        </Search>

    </div>
</template>



<script>
import Search from '@/components/Search.vue';

export default {
    components: {
        Search
    },
    data() {
        return {
            name: this.$route.query.name || '',
            avatar: this.$route.query.avatar || '',
            groupID: this.$route.query.groupID || '',
            members: [],
            searchModalVisible: false,
            showNameModal: false,
            showAvatarModal: false,
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
    },
    mounted() {
        if (!sessionStorage.getItem('token')) {
            this.$router.push("/");
            return;
        }
        this.getMembers();
        this.intervalId = setInterval(async () => {
          clearInterval(this.intervalId);
          await this.getMembers();
          this.intervalId = setInterval(this.getMembers, 1000);
        }, 1000);
    },
    beforeUnmount() {
        if (this.intervalId) {
          clearInterval(this.intervalId);
        }
    }
}

</script>