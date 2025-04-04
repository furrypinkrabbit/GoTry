<template>
  <el-container>
    <el-main>
      <el-form :model="form" class="upload-form" @submit.prevent="handleVideoUpload">
        <el-form-item label="选择视频" label-width="80px">
          <el-input v-model="form.videoName" readonly placeholder="请选择视频">
            <template #suffix>
              <el-button @click="openFileDialog">选择文件</el-button>
            </template>
          </el-input>
          <input type="file" ref="fileInput" accept="video/*" @change="onFileChange" style="display: none;">
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit">上传视频</el-button>
        </el-form-item>
      </el-form>
      <el-progress :percentage="progress" v-if="uploading"></el-progress>
      <el-alert :title="errorMessage" type="error" v-if="errorMessage"></el-alert>
      <!-- 调整视频窗口大小并居中 -->
      <div class="video-container">
        <!-- 添加 preload、autoplay 和 loop 属性 -->
        <video ref="videoPreview" width="640" height="480" controls v-if="videoUrl" preload="auto" autoplay loop muted></video>
      </div>
      <!-- 添加不可编辑的文字框 -->
      <el-input v-model="videoSummary" placeholder="视频摘要" v-if="videoUrl" readonly></el-input>
      <el-button @click="prevVideo" v-if="videoUrls.length > 0">上一个</el-button>
      <el-button @click="nextVideo" v-if="videoUrls.length > 0">下一个</el-button>
    </el-main>
    <!-- 加载界面 -->
    <div v-if="loadingInstance" class="loading-mask">
      <div class="loading-content">
        <p>正在上传视频，请稍候...</p>
        <el-button @click="cancelUpload">取消</el-button>
      </div>
    </div>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { ElMessage } from 'element-plus';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const activeIndex = ref(route.name?.toString() || 'uploadVideo');

const fileInput = ref<HTMLInputElement | null>(null);
const videoPreview = ref<HTMLVideoElement | null>(null);
const form = ref({
  videoName: '',
});
const progress = ref(0);
const uploading = ref(false);
const errorMessage = ref('');
const videoUrl = ref('');
const videoSummary = ref('');
const videoUrls = ref<string[]>([]);
const textSummaries = ref<string[]>([]);
// 新增历史上传记录数组
const uploadHistory = ref<{ videoName: string; textSummary: string; uploadTime: string; filePairs?: FilePair[] }[]>(
  JSON.parse(localStorage.getItem('uploadHistory') || '[]')
);

// 新增加载实例状态
const loadingInstance = ref<boolean>(false);

interface FilePair {
  video_url: string;
  text: string;
}

let currentIndex = 0;

const openFileDialog = () => {
  if (fileInput.value) {
    fileInput.value.click();
  }
};

const onFileChange = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (file) {
    form.value.videoName = file.name;
    const reader = new FileReader();
    reader.onload = function (e) {
      if (videoPreview.value) {
        videoUrl.value = e.target?.result as string;
        videoPreview.value.src = videoUrl.value;
        // 手动加载并播放视频
        videoPreview.value.load();
        videoPreview.value.play();
      }
    };
    reader.readAsDataURL(file);
  }
};

const handleVideoUpload = async () => {
  const file = (fileInput.value as HTMLInputElement).files?.[0];
  if (!file) {
    ElMessage.error('请选择视频文件');
    return;
  }

  uploading.value = true;
  errorMessage.value = '';

  // 显示加载动画
  loadingInstance.value = true;

  const formData = new FormData();
  formData.append('video', file);
  formData.append('name', file.name);
  formData.append('size', file.size.toString());
  formData.append('title', file.name);

  try {
    const token = authStore.token;
    const response = await fetch('http://localhost:3000/api/uploadVideo', {
      method: 'POST',
      headers: {
        'Authorization': `${token}`,
      },
      body: formData,
    });

    if (response.ok) {
      const data = await response.json();
      const checkProcessedVideo = async () => {
        const filePairsResponse = await fetch('http://localhost:3000/api/getVideo');
        if (filePairsResponse.ok) {
          const filePairsData = await filePairsResponse.json();
          const filePairs: FilePair[] = filePairsData.file_pairs;
          videoUrls.value = filePairs.map((pair: FilePair) => 'http://localhost:3000' + pair.video_url);
          textSummaries.value = filePairs.map((pair) => pair.text);

          currentIndex = 0;
          if (videoUrls.value.length > 0) {
            videoUrl.value = videoUrls.value[currentIndex];
            videoSummary.value = textSummaries.value[currentIndex];
            if (videoPreview.value) {
              videoPreview.value.src = videoUrl.value;
              // 手动加载并播放视频
              videoPreview.value.load();
              videoPreview.value.play();
            }
          }
          uploading.value = false;
          progress.value = 0;

          // 记录上传历史
          const uploadTime = new Date().toLocaleString();

          uploadHistory.value.push({
            videoName: file.name,
            textSummary: textSummaries.value.join(', '), // Combine text summaries into a single string
            filePairs: filePairs, // 新增存储多个视频文字对
            uploadTime,
          });

          // 将上传历史保存到本地存储
          localStorage.setItem('uploadHistory', JSON.stringify(uploadHistory.value));

          // 关闭加载动画
          loadingInstance.value = false;
        } else {
          setTimeout(checkProcessedVideo, 60000 * 3);
        }
      };
      checkProcessedVideo();
    } else {
      errorMessage.value = '视频上传失败';
      uploading.value = false;
      progress.value = 0;
      // 关闭加载动画
      loadingInstance.value = false;
    }
  } catch (error) {
    errorMessage.value = '上传视频时出错';
    uploading.value = false;
    progress.value = 0;
    // 关闭加载动画
    loadingInstance.value = false;
  }
};

const prevVideo = () => {
  if (videoUrls.value.length > 0) {
    currentIndex = (currentIndex - 1 + videoUrls.value.length) % videoUrls.value.length;
    videoUrl.value = videoUrls.value[currentIndex];
    videoSummary.value = textSummaries.value[currentIndex];
    if (videoPreview.value) {
      videoPreview.value.src = videoUrl.value;
      videoPreview.value.load();
      videoPreview.value.play();
    }
  }
};

const nextVideo = () => {
  if (videoUrls.value.length > 0) {
    currentIndex = (currentIndex + 1) % videoUrls.value.length;
    videoUrl.value = videoUrls.value[currentIndex];
    videoSummary.value = textSummaries.value[currentIndex];
    if (videoPreview.value) {
      videoPreview.value.src = videoUrl.value;
      videoPreview.value.load();
      videoPreview.value.play();
    }
  }
};

const handleSelect = (key: string) => {
  if (key === 'logout') {
    authStore.logout();
    router.push({ name: 'Home' });
  } else {
    router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) });
  }
};

// 取消上传功能
const cancelUpload = () => {
  loadingInstance.value = false;
  uploading.value = false;
  progress.value = 0;
  errorMessage.value = '上传已取消';
};

// 路由守卫，保存和恢复状态
const saveState = () => {
  localStorage.setItem('uploadState', JSON.stringify({
    uploading: uploading.value,
    progress: progress.value,
    errorMessage: errorMessage.value,
    videoUrl: videoUrl.value,
    videoSummary: videoSummary.value,
    videoUrls: videoUrls.value,
    textSummaries: textSummaries.value,
    currentIndex: currentIndex,
    loadingInstance: loadingInstance.value
  }));
};

const restoreState = () => {
  const state = localStorage.getItem('uploadState');
  if (state) {
    const { uploading: up, progress: prog, errorMessage: err, videoUrl: url, videoSummary: summary, videoUrls: urls, textSummaries: texts, currentIndex: idx, loadingInstance: load } = JSON.parse(state);
    uploading.value = up;
    progress.value = prog;
    errorMessage.value = err;
    videoUrl.value = url;
    videoSummary.value = summary;
    videoUrls.value = urls;
    textSummaries.value = texts;
    currentIndex = idx;
    loadingInstance.value = load;
  }
};

onBeforeRouteLeave(() => {
  saveState();
});

onBeforeRouteUpdate(() => {
  saveState();
});

onMounted(() => {
  restoreState();
});

</script>

<style scoped>
.upload-form {
  width: 100%;
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  background-color: #f5f5f5;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 视频容器居中 */
.video-container {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.loading-mask {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
}

.loading-content {
  background-color: white;
  padding: 20px;
  border-radius: 4px;
  text-align: center;
}
</style>
