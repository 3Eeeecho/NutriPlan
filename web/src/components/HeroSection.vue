<template>
  <section class="hero">
    <nav class="hero__nav">
      <div class="hero__nav-tag">NUTRIPLAN FOR LIFESTYLE</div>
      <div class="hero__nav-right">
        <button class="hero__btn hero__btn--nav hero__btn--soft" @click="router.push('/login')">登录</button>
        <button class="hero__btn hero__btn--nav hero__btn--outline" @click="router.push('/register')">免费注册</button>
      </div>
    </nav>

    <main class="hero__content">
      <section class="hero__copy">
        <h1 class="hero__title">
          健康饮食，<br />一拍即合。
        </h1>

        <p class="hero__subtitle">
          告别繁琐的手动查阅。随手一拍，智能解析餐盘中的热量与营养，帮你更快建立可持续的饮食节奏。
        </p>

        <div class="hero__feature-list">
          <article class="hero__feature-item">
            <h3>拍照识别</h3>
            <p>3 秒获取热量、蛋白质、脂肪与碳水估算。</p>
          </article>
          <article class="hero__feature-item">
            <h3>智能建议</h3>
            <p>根据体重目标与历史记录，动态调整每日摄入方案。</p>
          </article>
          <article class="hero__feature-item">
            <h3>一键规划</h3>
            <p>自动生成食谱、采购清单与每周复盘报告。</p>
          </article>
        </div>

        <div class="hero__actions">
          <button class="hero__btn hero__btn--cta hero__btn--primary" @click="router.push('/login')">开启智能饮食</button>
          <button ref="demoTriggerRef" class="hero__btn hero__btn--cta hero__btn--outline" @click="openDemoVideo">查看演示</button>
        </div>
      </section>

      <aside class="hero__visual" aria-label="产品展示图片">
        <figure class="hero__image-wrap">
          <img :src="featureImageMain" alt="NutriPlan 饮食与健康示意图" />
        </figure>

        <div class="hero__floating-card">
          <p class="hero__floating-label">今日摘要</p>
          <p class="hero__floating-value">热量达成率 82%</p>
          <p class="hero__floating-desc">连续 6 天完成蛋白质目标</p>
        </div>
      </aside>
    </main>

    <div v-if="isDemoVisible" class="hero__video-overlay" @click.self="closeDemoVideo">
      <div class="hero__video-modal" role="dialog" aria-modal="true" aria-label="产品演示视频">
        <div class="hero__video-header">
          <p class="hero__video-title">NutriPlan 产品演示</p>
          <button class="hero__video-close" @click="closeDemoVideo">关闭</button>
        </div>
        <video ref="demoVideoRef" class="hero__video-player" controls autoplay muted playsinline :poster="demoPoster">
          <source :src="demoVideoSrc" type="video/mp4" />
          当前浏览器不支持视频播放。
        </video>
      </div>
    </div>
  </section>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isDemoVisible = ref(false)
const demoVideoRef = ref(null)
const demoTriggerRef = ref(null)
const previousBodyOverflow = ref('')

const featureImageMain =
  'https://images.unsplash.com/photo-1498837167922-ddd27525d352?auto=format&fit=crop&w=1500&q=80'

const demoVideoSrc = '/videos/demo-placeholder.mp4'
const demoPoster = 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&w=1200&q=80'

const openDemoVideo = () => {
  isDemoVisible.value = true
}

const closeDemoVideo = () => {
  if (demoVideoRef.value) {
    demoVideoRef.value.pause()
    demoVideoRef.value.currentTime = 0
  }
  isDemoVisible.value = false

  nextTick(() => {
    if (demoTriggerRef.value) {
      demoTriggerRef.value.focus()
    }
  })
}

const onEscClose = (event) => {
  if (event.key === 'Escape' && isDemoVisible.value) {
    closeDemoVideo()
  }
}

watch(isDemoVisible, (visible) => {
  if (visible) {
    previousBodyOverflow.value = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return
  }

  document.body.style.overflow = previousBodyOverflow.value
})

onMounted(() => {
  window.addEventListener('keydown', onEscClose)
})

onBeforeUnmount(() => {
  document.body.style.overflow = previousBodyOverflow.value
  window.removeEventListener('keydown', onEscClose)
})
</script>

<style scoped>
.hero {
  --btn-radius: 12px;
  --btn-transition: transform 0.18s ease, box-shadow 0.18s ease, background-color 0.18s ease, border-color 0.18s ease;
  --btn-font-size: 14px;
  --btn-font-weight: 600;
  --btn-letter-spacing: 0.02em;
  position: fixed;
  inset: 0;
  z-index: 50;
  min-height: 100vh;
  min-height: 100dvh;
  overflow-y: auto;
  background: #d7dfd2;
  color: #111111;
}

.hero__nav {
  max-width: 1320px;
  height: 90px;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hero__nav-tag {
  font-size: 14px;
  letter-spacing: 0.02em;
  font-weight: 500;
}

.hero__nav-right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.hero__btn {
  border: 0;
  background: transparent;
  min-height: 38px;
  padding: 0 14px;
  border-radius: var(--btn-radius);
  font-size: var(--btn-font-size);
  font-family: inherit;
  letter-spacing: var(--btn-letter-spacing);
  font-weight: var(--btn-font-weight);
  line-height: 1;
  cursor: pointer;
  transition: var(--btn-transition);
}

.hero__btn--nav {
  min-height: 38px;
  padding: 0 14px;
  font-size: 13px;
  font-weight: 500;
}

.hero__btn--cta {
  min-height: 50px;
  padding: 0 22px;
}

.hero__btn--soft {
  color: #161616;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(22, 22, 22, 0.2);
}

.hero__btn--outline {
  border: 1px solid #161616;
  color: #161616;
  background: rgba(255, 255, 255, 0.28);
}

.hero__btn--primary {
  border: 1px solid #1f8f45;
  background: #2fa45b;
  color: #fff;
  box-shadow: 0 8px 20px rgba(47, 164, 91, 0.28);
}

.hero__btn--nav.hero__btn--outline {
  font-weight: 600;
  box-shadow: 0 4px 14px rgba(17, 17, 17, 0.08);
}

.hero__content {
  max-width: 1320px;
  margin: 0 auto;
  padding: 28px 32px 48px;
  min-height: calc(100dvh - 90px);
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 56px;
  align-items: center;
}

.hero__copy {
  max-width: 640px;
  animation: fadeInUp 0.7s ease both;
}

.hero__title {
  margin: 0;
  font-size: clamp(40px, 5.8vw, 84px);
  line-height: 1.02;
  letter-spacing: -0.03em;
  font-weight: 600;
}

.hero__subtitle {
  margin: 22px 0 0;
  font-size: clamp(18px, 1.8vw, 27px);
  line-height: 1.45;
  color: rgba(17, 17, 17, 0.9);
}

.hero__feature-list {
  margin-top: 26px;
  display: grid;
  gap: 14px;
}

.hero__feature-item h3 {
  margin: 0;
  font-size: 28px;
  line-height: 1.1;
  text-decoration: underline;
  text-underline-offset: 4px;
  font-weight: 600;
}

.hero__feature-item p {
  margin: 6px 0 0;
  font-size: 18px;
  line-height: 1.45;
  color: rgba(17, 17, 17, 0.9);
}

.hero__actions {
  margin-top: 30px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.hero__visual {
  position: relative;
  display: grid;
  place-items: center;
}

.hero__image-wrap {
  width: min(100%, 620px);
  aspect-ratio: 4 / 5;
  margin: 0;
  border-radius: 26px;
  overflow: hidden;
  box-shadow: 0 18px 42px rgba(0, 0, 0, 0.18);
}

.hero__image-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.hero__floating-card {
  position: absolute;
  left: 10%;
  bottom: 4%;
  width: min(360px, 78%);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(6px);
  padding: 16px 18px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.hero__floating-label {
  margin: 0;
  font-size: 12px;
  color: rgba(17, 17, 17, 0.7);
}

.hero__floating-value {
  margin: 6px 0 0;
  font-size: 20px;
  font-weight: 700;
}

.hero__floating-desc {
  margin: 6px 0 0;
  font-size: 14px;
  color: rgba(17, 17, 17, 0.72);
}

.hero__btn:hover {
  transform: translateY(-1px);
}

.hero__btn--soft:hover,
.hero__btn--outline:hover {
  background: rgba(255, 255, 255, 0.42);
}

.hero__btn--soft:hover {
  box-shadow: 0 4px 12px rgba(17, 17, 17, 0.1);
}

.hero__btn--primary:hover {
  background: #289150;
  box-shadow: 0 10px 22px rgba(40, 145, 80, 0.32);
}

.hero__btn:focus-visible {
  outline: 2px solid rgba(17, 17, 17, 0.5);
  outline-offset: 2px;
}

.hero__video-overlay {
  position: fixed;
  inset: 0;
  z-index: 90;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(0, 0, 0, 0.55);
  backdrop-filter: blur(3px);
}

.hero__video-modal {
  width: min(920px, 100%);
  border-radius: 16px;
  overflow: hidden;
  background: #0f1115;
  box-shadow: 0 22px 50px rgba(0, 0, 0, 0.35);
}

.hero__video-header {
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #171b22;
}

.hero__video-title {
  margin: 0;
  color: #f5f8ff;
  font-size: 14px;
  font-weight: 600;
}

.hero__video-close {
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: transparent;
  color: #f5f8ff;
  min-height: 30px;
  padding: 0 12px;
  border-radius: 8px;
  cursor: pointer;
}

.hero__video-player {
  display: block;
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #000;
}

@media (max-width: 1100px) {
  .hero__content {
    grid-template-columns: 1fr;
    gap: 32px;
    padding-top: 12px;
  }

  .hero__copy {
    max-width: 100%;
  }

  .hero__visual {
    width: min(100%, 700px);
  }
}

@media (max-width: 720px) {
  .hero__nav {
    height: 74px;
    padding: 0 16px;
  }

  .hero__nav-tag {
    font-size: 12px;
  }

  .hero__btn--nav {
    min-height: 34px;
    padding: 0 10px;
    font-size: 12px;
  }

  .hero__content {
    min-height: calc(100dvh - 74px);
    padding: 8px 16px 30px;
    gap: 22px;
  }

  .hero__feature-item h3 {
    font-size: 24px;
  }

  .hero__feature-item p {
    font-size: 16px;
  }

  .hero__floating-card {
    left: 6%;
    bottom: 4%;
    width: 88%;
  }
}

@media (max-width: 520px) {
  .hero__title {
    font-size: 44px;
  }

  .hero__subtitle {
    font-size: 18px;
  }

  .hero__actions {
    display: grid;
  }

  .hero__btn--cta {
    width: min(100%, 280px);
  }

  .hero__floating-card {
    position: static;
    width: 100%;
    margin-top: 14px;
  }

  .hero__video-overlay {
    padding: 12px;
  }

  .hero__video-title {
    font-size: 13px;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
