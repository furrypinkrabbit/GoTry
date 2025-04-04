<template>
    <el-container>
        <el-main>
            <el-alert :title="errorMessage" type="error" v-if="errorMessage"></el-alert>
            <video ref="videoPreview" width="320" height="240" controls v-if="videoUrl"></video>
            <el-input v-model="receivedText" type="textarea" placeholder="接收到的文本信息" style="margin-top: 20px;"></el-input>
            <el-button @click="fetchVideo">获取处理后的视频</el-button>
        </el-main>
    </el-container>
  </template>
  
  <script setup lang="ts">
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useAuthStore } from '../store/auth';
import axios from 'axios';

const videoPreview = ref<HTMLVideoElement | null>(null);
const errorMessage = ref('');
const videoUrl = ref('');
const receivedText = ref('');
const md5Hash = ref('');
const authStore = useAuthStore();
const fetchVideo = async () => {
    errorMessage.value = '';
    try {
        // 假设从本地存储或其他地方获取令牌
    const token = authStore.token; 
        
        const response = await axios.get('http:localhost:3000/api/getVideo', {
            headers: {
                Authorization: `${token}`
            }
        });
        // 处理响应
        if (response.data && response.data.file_pairs && response.data.file_pairs.length > 0) {
            const filePair = response.data.file_pairs[0];
            const videoPath = filePair.VideoURL;
            receivedText.value = filePair.Text;
            md5Hash.value = filePair.MD5Hash;
            videoUrl.value = videoPath;
            if (videoPreview.value) {
                videoPreview.value.src = videoPath;
            }
            ElMessage.success('视频获取成功');
        } else {
            errorMessage.value = '未找到匹配的视频和文字对';
        }
    } catch (error) {
        errorMessage.value = '获取视频时出错';
        console.error('获取视频出错:', error);
    }
};
</script>
  
  <style scoped>
  /* 可以根据需要调整样式 */
  video {
    margin-top: 20px;
  }
  </style>
  