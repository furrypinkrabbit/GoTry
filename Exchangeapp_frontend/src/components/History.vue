<template>
  <el-container>
    <el-main>
      <h1 class="centered-title">历史记录</h1>
      <div class="upload-form">
        <!-- 添加清空按钮 -->
        <el-button @click="clearHistory">清空历史记录</el-button>
        <el-table :data="uploadHistory" stripe>
          <el-table-column prop="videoName" label="视频名称"></el-table-column>
          <el-table-column label="视频描述与链接">
            <template #default="{ row }">
              <div v-for="(pair, index) in row.filePairs" :key="index">
                <a :href="'http://localhost:3000' + pair.video_url" target="_blank">{{ '点击查看视频' }}</a>
                <span>{{ pair.text }}</span>
                <span>{{ row.uploadTime }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="uploadTime" label="上传时间"></el-table-column>
        </el-table>
      </div>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '../store/auth';

const authStore = useAuthStore();
// 从本地存储中获取上传历史记录
const uploadHistory = ref<{ 
  videoName: string; 
  filePairs: { video_url: string; text: string }[]; 
  uploadTime: string 
}[]>(
  JSON.parse(localStorage.getItem('uploadHistory') || '[]')
);

// 清空历史记录的方法
const clearHistory = () => {
  // 清空本地存储中的上传历史记录
  localStorage.removeItem('uploadHistory');
  // 更新表格数据
  uploadHistory.value = [];
};

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

.centered-title {
  text-align: center;
}
</style>
