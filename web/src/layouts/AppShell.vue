<template>
  <div class="app-shell">
    <aside class="app-shell-sidebar">
      <div class="shell-profile">
        <div class="shell-avatar">{{ userAvatarText }}</div>
        <div class="shell-user">
          <div class="shell-name">{{ userDisplayName }}</div>
          <div class="shell-sub">NutriPlan 用户</div>
        </div>
      </div>

      <nav class="shell-menu">
        <button
          v-for="item in menuItems"
          :key="item.key"
          class="shell-item"
          :class="{ active: isActive(item) }"
          @click="router.push(item.path)"
        >
          <n-icon><component :is="item.icon" /></n-icon>
          <span>{{ item.label }}</span>
        </button>
      </nav>

      <button class="shell-logout" @click="handleLogout">
        <n-icon><LogOutOutline /></n-icon>
        <span>退出登录</span>
      </button>
    </aside>

    <main class="app-shell-main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import {
  PersonOutline,
  StarOutline,
  RestaurantOutline,
  BarChartOutline,
  BasketOutline,
  LogOutOutline
} from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const userDisplayName = computed(() => authStore.user?.username || "用户");
const userAvatarText = computed(() => (userDisplayName.value || "用").slice(0, 1).toUpperCase());

const menuItems = [
  { key: "profile", label: "个人档案", path: "/profile/view", icon: PersonOutline },
  { key: "favorites", label: "我的收藏", path: "/favorites", icon: StarOutline },
  { key: "home", label: "饮食记录", path: "/home", icon: RestaurantOutline },
  { key: "weekly", label: "周报告", path: "/weekly", icon: BarChartOutline },
  { key: "shopping", label: "购物清单", path: "/shopping", icon: BasketOutline }
];

const isActive = (item) => route.path === item.path || route.path.startsWith(item.path + "/");

const handleLogout = () => {
  authStore.logout();
  router.push("/login");
};
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: #f8fafc;
  display: grid;
  grid-template-columns: 250px minmax(0, 1fr);
  gap: 1.5rem;
  padding: 1rem;
}

.app-shell-sidebar {
  position: sticky;
  top: 1rem;
  align-self: start;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 1rem;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.shell-profile {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f1f5f9;
}

.shell-avatar {
  width: 42px;
  height: 42px;
  border-radius: 999px;
  background: #10b981;
  color: #ffffff;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.shell-name {
  font-weight: 700;
  color: #1e293b;
}

.shell-sub {
  color: #94a3b8;
  font-size: 0.78rem;
  margin-top: 0.2rem;
}

.shell-menu {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}

.shell-item,
.shell-logout {
  border: none;
  background: transparent;
  color: #334155;
  border-radius: 0.65rem;
  padding: 0.55rem 0.7rem;
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
  font-size: 0.92rem;
}

.shell-item:hover,
.shell-logout:hover {
  background: #f8fafc;
}

.shell-item.active {
  background: #ecfdf5;
  color: #047857;
  font-weight: 600;
}

.shell-logout {
  margin-top: auto;
  color: #b45309;
}

.app-shell-main {
  min-width: 0;
}

@media (max-width: 900px) {
  .app-shell {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .app-shell-sidebar {
    position: static;
  }
}
</style>
