<template>
    <div class="container mt-5" v-if="notBanned">
        <div class="center-container">
            <div class="card photo-card">
                <button v-if="isMe" @click="deletePhoto" class="btn btn-danger delete-button mb-2">
                    Delete Photo <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#trash-2" />
                    </svg>
                </button>

                <!-- <img :src="imgSrc" alt="Photo" class="card-img-top" /> -->
                <div class="card-body photo-details">
                    <div class="author">{{ authorName }}, {{ formattedDate }}</div>
                    <div class="card-text text-center bg-light fs-5">{{ caption }}</div>
                    <div class="actions">
                        <button @click="likePhoto" class="btn btn-sm btn-outline-primary ms-3">
                            {{ Liked ? 'Unlike' : 'Like' }}
                        </button>
                        <span class="like-counter">{{ LikeCount }} Likes <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#thumbs-up" />
                            </svg></span>
                        <button @click="commentPhoto" class="btn btn-sm btn-outline-secondary" data-bs-toggle="modal"
                            :data-bs-target="'#usersModal' + modalId">
                            Comment <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#message-circle" />
                            </svg>
                        </button>
                        <button @click="viewComments" class="btn btn-sm btn-outline-secondary">
                            View Comments <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#message-square" />
                            </svg>
                        </button>
                        <!-- <CommentModal :photoId="this.modalId" /> -->
                    </div>
                    <div class="comments">
                        <div v-if="showComments" v-for="comment in photoComments" :key="comment.commentId" class="comment">
                                <div class="comment-author">{{ comment.userId }}</div>
                                <div class="comment-text">{{ comment.text }}</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
// import CommentModal from '@/components/CommentModal.vue';

const token = sessionStorage.getItem('authToken');
export default {
    components: {
        // CommentModal,
    },
    props: {
        photoId: Number,
        authorName: String,
        caption: String,
        date: Number,
        likes: Number,
        isLiked: Boolean,
    },
    data() {
        return {
            authorId: 0,
            isMe: false,
            imgSrc: null,
            notBanned: true,
            Liked: this.isLiked,
            LikeCount: this.likes,
            modalId: String(this.photoId),
            showComments: false,
            photoComments: [],
        };
    },
    async mounted() {

        if (this.photoId) {
            try {
                // const response = await this.$axios.get(`/photos/${this.photoId}`, {
                //   headers: {
                //     Authorization: `Bearer ${token}`,
                //   },
                //   responseType: 'blob',
                // });
                // const imageUrl = URL.createObjectURL(response.data);
                // this.imgSrc = imageUrl;
                this.findAuthorId();
            } catch (error) {
                if (error.response) {
                    const statusCode = error.response.status;
                    this.notBanned = false;
                    switch (statusCode) {
                        case 401:
                            console.error('Unauthorized:', error.response.data);
                            break;
                        case 403:
                            console.error('Forbidden Action:', error.response.data);
                            break;
                        case 404:
                            console.error('Not Found:', error.response.data);
                            break;
                        default:
                            console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                    }
                } else {
                    console.error('Error:', error);
                }
            }
        }
    },
    computed: {
        formattedDate() {
            const date = new Date(this.date);
            return date.toLocaleString();
        },
    },
    methods: {
        async findAuthorId() {
            try {
                const userId = this.$route.params.userId;
                const hasStreamSegment = this.$route.path.includes('/stream');
                if (userId == token && !hasStreamSegment) {
                    this.isMe = true;
                };
            }
            catch (error) {
                console.error(error, "Error during the author id retrival.")
            }
        },
        async deletePhoto() {
            try {
                await this.$axios.delete(`/posts/${this.photoId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }
                },);
                location.reload();
            }
            catch (error) {
                console.error(error, "Unable to delete the post")
            }
        },
        async likePhoto() {
            try {
                if (!this.Liked) {
                    await this.$axios.post(`/posts/${this.photoId}/likes/self`, {
                    }, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                } else {
                    await this.$axios.delete(`/posts/${this.photoId}/likes/self`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                }
                this.LikeCount = this.Liked ? this.LikeCount - 1 : this.LikeCount + 1;
                this.Liked = !this.Liked;
            } catch (error) {
                console.error(error, "Error during the likes operation.")
            }
        },
        async viewComments() {
            try {
                this.showComments = !this.showComments;
                if (this.showComments) {
                    const response = await this.$axios.get(`/posts/${this.photoId}/comments/`, {
                        headers: {
                            Authorization: `Bearer ${token}`
                        }
                    });
                    this.photoComments = response.data.comments !== null ? response.data.comments : [];
                }
            } catch (error) {
                console.error(error, "Error during the comments operation.")
            }
        }
    },
};
</script>

<style scoped>
.center-container {
    display: flex;
    justify-content: center;
    align-items: center;
}

.photo-card {
    border: 3px solid #6d6969;
    border-radius: 4px;
    padding: 10px;
    width: 500px;
    text-align: center;
    font-family: 'Arial', sans-serif;
}

.photo-details {
    margin-top: 10px;
}

.author {
    font-size: 20px;
    margin-bottom: 5px;
}

.actions{
    display: flex;
    justify-content: space-between;
    margin: 15px;
}

.comments {
    display: flex;
    flex-direction: column;
    gap: 10px; 
    max-width: 400px; 
    margin: 0 auto; 
}

.comment {
    display: flex;
    flex-direction: column;
    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 10px;
    background-color: #f9f9f9;
    text-align: right;
}

.comment-author {
    font-weight: bold;
    margin-bottom: 5px;
    color: #333;
}

.comment-text {
    color: #555;
    line-height: 1.5;
}

.like-counter {
    margin-left: 2px;
    border: 2px solid #d102027a;
    border-radius: 4px;
    padding: 8px;
}

.caption {
    display: flex;
    align-items: center;
    margin-top: 10px;
}

.caption-border {
    flex: 1;
    height: 3px;
    background-color: #1a1212;
    padding: 4px;
    margin-top: 10px;
    margin-bottom: 10px;

}

.caption-text {
    padding: 0 10px;
}
</style>