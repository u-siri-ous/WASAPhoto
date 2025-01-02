<template>
    <div class="modal fade" :id="'userListModal' + postId" tabindex="-1" aria-labelledby="ModalLabel" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">{{ typeOfList }}</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <ul class="user-list">
                        <li v-for="user in users" :key="user.userId" class="user-item">
                            <span>{{ user.username }}</span>
                        </li>
                    </ul>
                </div>

            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        users: Array,
        postId: Number,
        typeOfList: String,
    },
};
</script>

<style scoped></style>