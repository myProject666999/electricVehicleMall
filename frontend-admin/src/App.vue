<script setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()

// 不需要布局的页面（登录）
const noLayoutPages = ['/login']

const showLayout = computed(() => {
  return !noLayoutPages.includes(route.path)
})
</script>

<template>
  <!-- 有布局的页面 -->
  <template v-if="showLayout">
    <router-view />
  </template>

  <!-- 无布局的页面（登录） -->
  <template v-else>
    <router-view />
  </template>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  background-color: #f5f7fa;
}
</style>
