<template>
    <div class="container">
        <h2>Stream</h2>
        <div class="photo-grid">
            <PhotoComponent v-for="photo in photoList" :key="photo.photoId" :photoId="photo.photoId"
                :date="photo.timeOfCreation" :authorName="photo.author" :authorId="photo.authorId" :likes="photo.numberOfLikes"
                :comments="photo.numberOfComments" :caption="photo.caption" :isLiked="photo.isliked" />
        </div>
    </div>
    <UserList :users="UserList" :postId="token" :typeOfList="TypeOfList" />
</template>

<script>
import PhotoComponent from '../components/PhotoComponent.vue';
import UserList from '../components/UserList.vue';

const token = sessionStorage.getItem('authToken');
export default {
    data() {
        return {
            photoList: [],
            UserList: {},
            TypeOfList: '',
        };
    },
    async mounted() {
        try {
            const userId = this.$route.params.userId;
            const response = await this.$axios.get(`/users/${userId}/stream?followed=true`, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            this.photoList = response.data.stream;
        } catch (error) {
            if (error.response) {
                const statusCode = error.response.status;
                switch (statusCode) {
                    case 401:
                        console.error('Unauthorized:', error.response.data);
                        break;
                    default:
                        console.error('Error:', error.response.data);
                        break;
                }
            } else {
                console.error('Error:', error.message);
            }
        }
    },
    methods: {
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
        }
    },
    components: {
        PhotoComponent,
        UserList,
    },
};
</script>