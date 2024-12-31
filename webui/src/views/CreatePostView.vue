<template>
    <div class="container mt-5">
        <h2 class="text-center mb-4">Upload Your Photo</h2>
        <form @submit.prevent="uploadPhoto" class="upload-form">
            <div class="mb-3">
                <label for="photo" class="form-label">Select Photo (JPG only)</label>
                <input type="file" id="photo" class="form-control" accept=".jpg" @change="onFileChange" required />
                <div class="invalid-feedback" v-if="!photo">Please select a photo before uploading.</div>
            </div>

            <div class="mb-3">
                <label for="caption" class="form-label">Add a Caption</label>
                <textarea id="caption" class="form-control" v-model="caption"
                    placeholder="Write something about your photo..."></textarea>
            </div>

            <div class="d-flex justify-content-center mt-4">
                <button type="submit" class="btn btn-primary" :disabled="!photo">
                    Upload Photo
                </button>
            </div>

            <div class="mt-4">
                <div v-if="uploadSuccess" class="alert alert-success text-center">
                    {{ photoUploadResult }}
                </div>
                <div v-else-if="photoUploadResult && !uploadSuccess" class="alert alert-danger text-center">
                    {{ photoUploadResult }}
                </div>
            </div>
        </form>
    </div>
</template>

<script>
const token = sessionStorage.getItem('authToken');

export default {
    data() {
        return {
            photo: null,
            caption: '',
            uploadSuccess: false,
            photoUploadResult: '',
        };
    },
    methods: {
        onFileChange(event) {
            this.photo = event.target.files[0];
        },
        async uploadPhoto() {
            if (!this.photo) {
                this.photoUploadResult = "Please select a photo to upload.";
                this.uploadSuccess = false;
                return;
            }

            const formData = new FormData();
            formData.append('photo', this.photo);
            formData.append('caption', this.caption);

            try {
                const response = await this.$axios.post(`/posts/`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        'Authorization': `Bearer ${token}`,
                    },
                });

                this.photoUploadResult = "Your photo has been uploaded successfully!";
                this.uploadSuccess = true;
                this.photo = null;
                this.caption = '';
            } catch (error) {
                this.uploadSuccess = false;
                if (error.response) {
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 401:
                            this.photoUploadResult = "You need to log in to upload photos.";
                            break;
                        default:
                            this.photoUploadResult = "An unexpected error occurred. Please try again.";
                    }
                } else {
                    this.photoUploadResult = "Unable to connect to the server. Check your internet connection.";
                }
            }
        },
    },
};
</script>