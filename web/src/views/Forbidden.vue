<template>
  <main class="forbidden-page">
    <section class="forbidden-panel">
      <p class="eyebrow">403</p>
      <h1>权限不足</h1>
      <p>当前账号没有访问该页面的权限。</p>
      <button type="button" @click="goHome">返回首页</button>
    </section>
  </main>
</template>

<script setup>
import { useRouter } from "vue-router";

const router = useRouter();

const goHome = () => {
  router.push("/home");
};
</script>

<style scoped>
.forbidden-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #eef5f1;
  padding: 24px;
}

.forbidden-panel {
  width: min(420px, 100%);
  background: #ffffff;
  border: 1px solid rgba(31, 86, 68, 0.12);
  border-radius: 8px;
  padding: 32px;
  box-shadow: 0 16px 32px rgba(31, 86, 68, 0.08);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7c63;
  font-weight: 700;
  letter-spacing: 0;
}

h1 {
  margin: 0;
  color: #1f5445;
  font-size: 28px;
}

p {
  color: #587a6f;
  line-height: 1.7;
}

button {
  border: none;
  border-radius: 6px;
  background: #2f7c63;
  color: #ffffff;
  padding: 10px 16px;
  font-weight: 700;
  cursor: pointer;
}
</style>
