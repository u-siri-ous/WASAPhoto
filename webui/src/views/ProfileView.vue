<template>
    <div class="container mt-5">
        <h1 class="display-4" style="font-size: 50px;">{{ userName }}</h1>
        <div v-if="found">
            <div>
                <div v-if="!isItMe">
                    <div class="btn-group mt-1">
                        <button @click="toggleFollow" class="btn btn-warning">
                            {{ isFollowed ? 'Unfollow' : 'Follow' }} <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#user-plus" />
                            </svg>
                        </button>
                        <button @click="toggleBan" class="btn btn-danger">
                            {{ isBanned ? 'Unban' : 'Ban' }} <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#slash" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
            <div style="font-size: 20px;" class="container mt-2">
                <div class="row bg-light p-4 shadow-lg">
                    <div class="row">Followers: {{ followCount }}</div>
                    <div class="row">Followed: {{ followedCount }}</div>
                    <div class="row">Photos: {{ photoCount }}</div>
                </div>
            </div>

        </div>
        <hr />
        <!-- <div class="photos">
            <PhotoCard v-for="photo in photoList" :key="photo.photoId" :photoId="photo.photoId" :date="photo.creation"
                :authorName="photo.author.username" :likes="photo.numberOfLikes" :caption="photo.caption" :Liked="photo.isliked"/>
        </div> -->
    </div>
</template>
  
<script>
// import PhotoCard from '@/components/PhotoCard.vue';
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
            isBanned: false,
            isFollowed: false,
            isItMe: false,
            photoList: [],
            reloadFlag: true,
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
                const response = await this.$axios.get(`/users/${userId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });

                this.found = true;
                this.userName = response.data.username;
                this.followCount = response.data.numberOfFollowers;
                this.followedCount = response.data.accountsFollowed;
                this.photoCount = response.data.numberOfPosts;
                this.isBanned = response.data.isBlocked;
                this.isFollowed = response.data.isFollowed;

                const response1 = await this.$axios.get(`/users/${userId}/photos/`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                });

                this.photoList = response1.data.posts;
                console.log(this.photoList)

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
                            this.userName = "You have been banned by the user"
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
            // frontend
            this.isFollowed = !this.isFollowed;
            // backend
            const userId = this.$route.params.userId;
            const token = sessionStorage.getItem('authToken');
            try {
                if (this.isFollowed) {
                    this.followCount += 1;
                    await this.$axios.put(`/followed_users/${userId}`, {
                    }, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                } else {
                    this.followCount -= 1;
                    await this.$axios.delete(`/followed_users/${userId}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                }
            } catch (error) {
                console.error(error, "Error modifying follow status.")
            }

        },
        async toggleBan() {
            // frontend
            this.isBanned = !this.isBanned;
            // backend
            const userId = this.$route.params.userId;
            const token = sessionStorage.getItem('authToken');
            try {
                if (this.isBanned) {
                    await this.$axios.put(`/banned_users/${userId}`, {
                    }, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                } else {
                    await this.$axios.delete(`/banned_users/${userId}`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });

                }
            } catch (error) {
                console.error(error, "Error modifying ban status.")
            }
        },
    },
    // components: {
    //     PhotoCard,
    // },
};
</script>
  
<style scoped>
.user-info {
    text-align: center;
    font-size: 20px;
}

hr {
    margin: 20px 0;
}

.photos {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
}

.photos .photo-card {
    margin-bottom: 30px;
}
</style>
  
  