<template>
    <div :class="['comment-line', { 'hover-enabled': this.isMe }]">
        <div class="comment w-100">
            <div class="comment-author font-weight-bold">{{ this.authorUsername }}</div>
            <div class="comment-text">{{ this.commentText }}</div>
        </div>
        <button @click="deleteComment" class="btn btn-danger delete-button h-100" :class="{ 'invisible': !this.isMe }">
            <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#trash" />
            </svg>
        </button>
    </div>
</template>


<script>
const token = sessionStorage.getItem('authToken');

export default {
    props: {
        photoId: Number,
        commentId: Number,
        authorUsername: String,
        authorId: Number,
        commentText: String,
    },
    data() {
        return {
            photoId: this.photoId,
            commentId: this.commentId,
            authorUsername: this.authorUsername,
            authorId: this.authorId,
            commentText: this.commentText,
            isMe: false,
        };
    },
    async mounted() {
        if (this.authorId == token) {
            this.isMe = true;
        }
    },
    methods: {
        async deleteComment() {
            try {
                await this.$axios.delete(`/posts/${this.photoId}/comments/${this.commentId}`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    }
                },);
                this.$parent.viewComments(true);
            }
            catch (error) {
                console.error(error, "Unable to delete the comment")
            }
        },
    }
}
</script>

<style scoped>
.comment-line {
    position: relative;
    display: flex;
    align-items: center;
}

.delete-button {
    position: absolute;
    right: 0;
    transform: translateX(100%);
    opacity: 0;
    transition: transform 0.3s ease, opacity 0.3s ease;
    pointer-events: none; 
}

.hover-enabled:hover .delete-button {
    transform: translateX(0);
    opacity: 1;
    pointer-events: auto; 
}

.delete-button.invisible {
    display: none; 
}


.comment {
    flex-direction: column;
    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 10px;
    background-color: #f9f9f9;
    text-align: left;
    flex: 1;
    transition: margin-right 0.3s ease;
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
</style>