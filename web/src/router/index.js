import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/store/auth";

const routes = [
  {
    path: "/",
    redirect: "/welcome",
  },
  {
    path: "/welcome",
    name: "Welcome",
    component: () => import("@/components/HeroSection.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("@/views/Register.vue"),
    meta: { requiresAuth: false },
  },
  {
    path: "/",
    component: () => import("@/layouts/AppShell.vue"),
    meta: { requiresAuth: true },
    children: [
      {
        path: "home",
        name: "Home",
        component: () => import("@/views/Home.vue"),
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/Profile.vue"),
      },
      {
        path: "profile/view",
        name: "ProfileView",
        component: () => import("@/views/ProfileView.vue"),
      },
      {
        path: "recipes",
        name: "RecipeRecommend",
        component: () => import("@/views/RecipeRecommend.vue"),
      },
      {
        path: "recipes/:id",
        name: "RecipeDetail",
        component: () => import("@/views/RecipeDetail.vue"),
      },
      {
        path: "favorites",
        name: "FavoriteList",
        component: () => import("@/views/FavoriteList.vue"),
      },
      {
        path: "intake",
        name: "IntakeRecord",
        component: () => import("@/views/IntakeRecord.vue"),
      },
      {
        path: "weekly-planner",
        name: "WeeklyPlanner",
        component: () => import("@/views/WeeklyPlanner.vue"),
      },
      {
        path: "weekly",
        name: "WeeklyReport",
        component: () => import("@/views/WeeklyReport.vue"),
      },
      {
        path: "shopping",
        name: "ShoppingList",
        component: () => import("@/views/ShoppingList.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.meta.requiresAuth;

  // 如果路由需要认证
  if (requiresAuth) {
    // 检查是否有 token
    if (authStore.isAuthenticated) {
      next();
    } else {
      // 没有 token，重定向到登录页
      next({ name: "Login", query: { redirect: to.fullPath } });
    }
  } else {
    // 如果已经登录，访问登录/注册页时重定向到首页
    if (
      authStore.isAuthenticated &&
      (to.name === "Login" || to.name === "Register")
    ) {
      next({ name: "Home" });
    } else {
      next();
    }
  }
});

export default router;
