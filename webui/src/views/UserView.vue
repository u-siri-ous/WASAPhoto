<template>
    <div class="profile-container">
        <h1 class="profile-username">{{ userName }}</h1>
        
        <div v-if="found" class="user-actions">
            <div class="action-buttons">
                <button v-if="!isItMe" @click="toggleFollow" class="btn col follow-btn">
                    {{ isFollowed ? 'Unfollow' : 'Follow' }}
                    <svg class="icon">
                        <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                    </svg>
                </button>
                <button v-if="!isItMe" @click="toggleBlock" class="btn col block-btn">
                    {{ isBlocked ? 'Unblock' : 'Block' }}
                    <svg class="icon">
                        <use href="/feather-sprite-v4.29.0.svg#slash" />
                    </svg>
                </button>
                <button v-if="isItMe" @click="viewBlock" class="btn col block-btn" data-bs-toggle="modal" :data-bs-target="'#userListModal' + token">
                    View Blocked
                    <svg class="icon">
                        <use href="/feather-sprite-v4.29.0.svg#user-x" />
                    </svg>
                </button>
            </div>

            <div class="user-info">
                <div class="info-card">
                    <button @click="viewFollowers" class="btn btn-stats btn-hover" data-bs-toggle="modal" :data-bs-target="'#userListModal' + token"><strong>Followers:</strong> {{ followCount }}</button>
                    <button @click="viewFollowing" class="btn btn-stats btn-hover" data-bs-toggle="modal" :data-bs-target="'#userListModal' + token"><strong>Followed:</strong> {{ followedCount }}</button>
                    <button class="btn btn-stats"><strong>Photos:</strong> {{ photoCount }}</button>
                </div>
            </div>
        </div>

        <hr />

        <div class="photo-gallery">
            <PhotoComponent 
                v-for="photo in photoList" 
                :key="photo.photoId" 
                :photoId="photo.photoId" 
                :date="photo.timeOfCreation" 
                :authorName="photo.author" 
                :likes="photo.numberOfLikes" 
                :comments="photo.numberOfComments"
                :caption="photo.caption" 
                :isLiked="photo.isliked" 
            />
        </div>
    </div>
    <UserList :users="UserList" :postId="token" :typeOfList="TypeOfList" />
</template>


<script>
import PhotoComponent from '../components/PhotoComponent.vue';
import UserList from '../components/UserList.vue';

const token = sessionStorage.getItem('authToken');

export default {
    mounted() {
        if (localStorage.getItem('reloadedstream')) {
            localStorage.removeItem('reloadedstream');
        } else {
            localStorage.setItem('reloadedstream', '1');
            location.reload();
        }
    },
    data() {
        return {
            userName: '',
            found: false,
            followCount: 0,
            followedCount: 0,
            photoCount: 0,
            isBlocked: false,
            isFollowed: false,
            isItMe: false,
            photoList: [],
            reloadFlag: true,
            UserList: {},
            TypeOfList: '',
        };
    },
    watch: {
        '$route.params.userId'(newParam, oldParam) {
            if (newParam !== oldParam) {
                this.refresh();
            }
        },
    },
    async mounted() {
        const token = sessionStorage.getItem('authToken');
        const userId = this.$route.params.userId;
        this.isItMe = (userId == token);
        this.fetchUserData();
    },
    methods: {
        refresh() {
            location.reload();
        },
        async fetchUserData() {
            const token = sessionStorage.getItem('authToken');
            const userId = this.$route.params.userId;
            try {
                const getUserResponse = await this.$axios.get(`/users/${userId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });

                this.found = true;
                this.userName = getUserResponse.data.username;
                this.followCount = getUserResponse.data.numberOfFollowers;
                this.followedCount = getUserResponse.data.accountsFollowed;
                this.photoCount = getUserResponse.data.numberOfPosts;
                this.isBlocked = getUserResponse.data.isBlocked;
                this.isFollowed = getUserResponse.data.isFollowed;

                const getStreamResponse = await this.$axios.get(`/users/${userId}/stream`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                this.photoList = getStreamResponse.data.stream;

            } catch (error) {
                if (error.response) {
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 400:
                            console.error('Bad request');
                            this.userName = "You have to login first"
                        case 401:
                            console.error('Access Unauthorized:', error.response.data);
                            this.userName = "You are not logged in"
                            break;
                        case 403:
                            console.error('Access Forbidden:', error.response.data);
                            this.userName = "You have been blocked by the user"
                            break;
                        case 404:
                            console.error('Not Found:', error.response.data);
                            if (userId === "null") {
                                this.userName = "You are not logged in";
                            }
                            else {
                                this.userName = "User not found";
                            }
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                    }
                } else {
                    console.error('Error:', error);
                }
            }
        },
        async toggleFollow() {
            this.isFollowed = !this.isFollowed;
            const userId = this.$route.params.userId;
            const token = sessionStorage.getItem('authToken');
            try {
                if (this.isFollowed) {
                    this.followCount += 1;
                    await this.$axios.put(`/following/${userId}`, {
                    }, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                } else {
                    this.followCount -= 1;
                    await this.$axios.delete(`/following/${userId}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                }
            } catch (error) {
                console.error(error, "Error modifying follow status.")
            }

        },
        async toggleBlock() {
            this.isBlocked = !this.isBlocked;
            const userId = this.$route.params.userId;
            const token = sessionStorage.getItem('authToken');
            try {
                if (this.isBlocked) {
                    await this.$axios.put(`/blocked/${userId}`, {
                    }, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                } else {
                    await this.$axios.delete(`/blocked/${userId}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });

                }
            } catch (error) {
                console.error(error, "Error during the block operation.")
            }
        },
        async viewLikes(photoId) {
            try {
                const response = await this.$axios.get(`/posts/${photoId}/likes/self`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                const users = response.data.users !== null ? response.data.users : [];
                this.UserList = users;
                this.TypeOfList = 'Likes';
            } catch (error) {
                console.error(error, "Error while showing the likes.")
            }
        },
        async viewFollowers() {
            try {
                const response = await this.$axios.get(`/followers/${this.$route.params.userId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                const users = response.data.users !== null ? response.data.users : [];
                this.UserList = users;
                this.TypeOfList = 'Followers';
            } catch (error) {
                console.error(error, "Error while showing the followers.")
            }
        },
        async viewFollowing() {
            try {
                const response = await this.$axios.get(`/following/${this.$route.params.userId}`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                const users = response.data.users !== null ? response.data.users : [];
                this.UserList = users;
                this.TypeOfList = 'Following';
            } catch (error) {
                console.error(error, "Error while showing the following.")
            }
        },
        async viewBlock() {
            try {
                const response = await this.$axios.get(`/blocked/`, {
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
                const users = response.data.users !== null ? response.data.users : [];
                this.UserList = users;
                this.TypeOfList = 'Blocked';
            } catch (error) {
                console.error(error, "Error while showing the blocked users.")
            }
        }
    },
    components: {
        PhotoComponent,
        UserList,
    },
};
</script>
  
<style scoped>
.profile-container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Roboto', sans-serif;
}

.profile-username {
    font-size: 2.5rem;
    font-weight: bold;
    text-align: center;
    margin-bottom: 20px;
    color: #333;
}

.user-actions {
    text-align: center;
    margin-bottom: 20px;
}

.action-buttons {
    display: flex;
    justify-content: center;
    gap: 10px;
    margin-bottom: 20px;
}

.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 6px 12px;
    font-size: 0.875rem;
    font-weight: bold;
    border: none;
    border-radius: 4px;
    transition: background-color 0.3s ease, transform 0.2s ease;
    cursor: pointer;
}

.follow-btn {
    background-color: #ffc107;
    color: #333;
}

.follow-btn:hover {
    background-color: #e0a800;
}

.block-btn {
    background-color: #dc3545;
    color: #fff;
}

.block-btn:hover {
    background-color: #c82333;
}

.user-info {
    background-color: #f8f9fa;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    text-align: center;
}

.info-card {
    display: flex;
    justify-content: space-between;
}

.info-card p {
    margin: 0;
    font-size: 1rem;
    color: #555;
}

.photo-gallery {
    display: block;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.icon {
    width: 16px;
    height: 16px;
    margin-left: 5px;
}

.btn-stats {
    background-color: #7a7a7a9e;
    color: #333;
}

.btn-stats.btn-hover:hover {
    background-color: #6767679e;
}
</style>
  
  