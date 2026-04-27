<template>
  <div class="app-shell" :class="{ expanded: 侧栏展开 }">
    <aside class="app-shell-sidebar" @mouseenter="侧栏展开 = true" @mouseleave="侧栏展开 = false">

      <div class="shell-profile">
        <div class="shell-avatar">{{ userAvatarText }}</div>
        <div class="shell-user">
          <div class="shell-name">{{ userDisplayName }}</div>
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
          <span class="label-text">{{ item.label }}</span>
        </button>
      </nav>

      <button class="shell-logout" @click="handleLogout">
        <n-icon><LogOutOutline /></n-icon>
        <span class="label-text">退出登录</span>
      </button>

      <button class="shell-record" @click="router.push('/intake')">
        <n-icon><AddCircleOutline /></n-icon>
        <span class="label-text">新建记录</span>
      </button>
    </aside>

    <main class="app-shell-main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { NIcon } from "naive-ui";
import {
  AddCircleOutline,
  PersonOutline,
  StarOutline,
  RestaurantOutline,
  BarChartOutline,
  BasketOutline,
  LogOutOutline,
  LeafOutline
} from "@vicons/ionicons5";
import { useAuthStore } from "@/store/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const 侧栏展开 = ref(false);

const userDisplayName = computed(() => authStore.user?.username || "用户");
const userAvatarText = computed(() => (userDisplayName.value || "用").slice(0, 1).toUpperCase());

const menuItems = computed(() => {
  const items = [
  { key: "profile", label: "个人档案", path: "/profile/view", icon: PersonOutline },
  { key: "favorites", label: "我的收藏", path: "/favorites", icon: StarOutline },
  { key: "home", label: "饮食记录", path: "/home", icon: RestaurantOutline },
  { key: "weekly", label: "周报告", path: "/weekly", icon: BarChartOutline },
  { key: "shopping", label: "购物清单", path: "/shopping", icon: BasketOutline }
  ];
  if (authStore.isAdmin) {
    items.push({ key: "admin", label: "管理后台", path: "/admin", icon: LeafOutline });
  }
  return items;
});

const isActive = (item) => route.path === item.path || route.path.startsWith(item.path + "/");

const handleLogout = () => {
  authStore.logout();
  router.push("/login");
};
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: #eaf2ee;
  display: grid;
  grid-template-columns: 84px minmax(0, 1fr);
  gap: 1rem;
  padding: 0.85rem;
}

.app-shell.expanded {
  grid-template-columns: 84px minmax(0, 1fr);
}

.app-shell-sidebar {
  position: sticky;
  top: 0.85rem;
  align-self: start;
  min-height: calc(100vh - 1.7rem);
  width: 68px;
  overflow: hidden;
  background: linear-gradient(180deg, rgba(247, 252, 249, 0.96), rgba(240, 248, 244, 0.94));
  border: 1px solid rgba(28, 86, 68, 0.12);
  border-radius: 1.15rem;
  padding: 1rem 0.7rem;
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
  box-shadow: 0 12px 28px rgba(31, 86, 68, 0.08);
  transition: width 0.22s ease, padding 0.22s ease;
  will-change: width;
}

.app-shell.expanded .app-shell-sidebar {
  width: 234px;
  padding: 1rem 0.9rem;
  z-index: 30;
}

.org-card {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  padding: 0.2rem 0.25rem 0.7rem;
  border-bottom: 1px solid rgba(26, 84, 66, 0.12);
}

.org-logo {
  width: 38px;
  height: 38px;
  border-radius: 999px;
  display: grid;
  place-items: center;
  background: #2f7c63;
  color: #ffffff;
}

.org-title {
  margin: 0;
  color: #214f42;
  font-weight: 700;
}

.org-subtitle {
  margin: 2px 0 0;
  color: #6f9587;
  font-size: 0.66rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.shell-profile {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 0.15rem 0.1rem 0.75rem;
  border-bottom: 1px solid rgba(26, 84, 66, 0.12);
}

.app-shell.expanded .shell-profile {
  justify-content: flex-start;
  gap: 0.75rem;
  padding: 0.15rem 0.25rem 0.75rem;
}

.shell-avatar {
  width: 42px;
  height: 42px;
  border-radius: 999px;
  background: linear-gradient(145deg, #2f7c63, #3d9a79);
  color: #ffffff;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.shell-name {
  font-weight: 700;
  color: #1f5445;
}

.shell-sub {
  color: #709486;
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
  color: #406f61;
  border-radius: 0.85rem;
  padding: 0.62rem;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s;
  font-size: 0.92rem;
}

.app-shell.expanded .shell-item,
.app-shell.expanded .shell-logout {
  padding: 0.62rem 0.72rem;
  justify-content: flex-start;
  gap: 0.5rem;
  text-align: left;
}

.shell-user {
  width: 0;
  opacity: 0;
  overflow: hidden;
  transition: width 0.22s ease, opacity 0.2s ease;
}

.app-shell.expanded .shell-user {
  width: 120px;
  opacity: 1;
}

.label-text {
  max-width: 0;
  opacity: 0;
  overflow: hidden;
  white-space: nowrap;
  transform: translateX(-4px);
  transition: max-width 0.24s ease, opacity 0.2s ease, transform 0.2s ease;
}

.app-shell.expanded .label-text {
  max-width: 130px;
  opacity: 1;
  transform: translateX(0);
}

.shell-item:hover,
.shell-logout:hover {
  background: #e7f2ec;
}

.shell-item.active {
  background: #ffffff;
  color: #225944;
  font-weight: 600;
  box-shadow: inset 0 0 0 1px rgba(35, 99, 79, 0.15), 0 8px 16px rgba(31, 86, 68, 0.09);
}

.shell-logout {
  margin-top: 0.4rem;
  color: #af5f2e;
}

.shell-record {
  margin-top: auto;
  border: none;
  border-radius: 999px;
  height: 48px;
  background: linear-gradient(140deg, #2f7c63, #3a9a78);
  color: #ffffff;
  font-size: 0.95rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(43, 117, 91, 0.28);
}

.app-shell.expanded .shell-record {
  gap: 8px;
}

.app-shell-main {
  min-width: 0;
}

@media (max-width: 900px) {
  .app-shell {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .app-shell.expanded {
    grid-template-columns: 1fr;
  }

  .app-shell-sidebar {
    position: static;
    width: auto;
    overflow: visible;
  }

  .shell-user,
  .label-text {
    width: auto;
    max-width: none;
    opacity: 1;
    overflow: visible;
    transform: none;
  }

  .shell-profile,
  .shell-item,
  .shell-logout,
  .shell-record {
    justify-content: flex-start;
    gap: 0.5rem;
    text-align: left;
  }
}
</style>
