const FONT_STACK = 'Segoe UI, Microsoft YaHei, PingFang SC, Helvetica, Arial, sans-serif'

const MEAL_META = {
  breakfast: {
    key: 'breakfast',
    label: '\u65e9\u9910',
    start: '#fff4d6',
    end: '#ffd39f',
    accent: '#c96f00',
  },
  lunch: {
    key: 'lunch',
    label: '\u5348\u9910',
    start: '#e6fff0',
    end: '#9fe3be',
    accent: '#0f8b5f',
  },
  dinner: {
    key: 'dinner',
    label: '\u665a\u9910',
    start: '#e7f0ff',
    end: '#a9c3ff',
    accent: '#3454c5',
  },
  snack: {
    key: 'snack',
    label: '\u52a0\u9910',
    start: '#fff1f4',
    end: '#ffb5c5',
    accent: '#c73f67',
  },
  default: {
    key: 'default',
    label: '\u9910\u98df',
    start: '#eef2f7',
    end: '#c8d2df',
    accent: '#4d647c',
  },
}

const CATEGORY_META = {
  fruit: {
    label: '\u6c34\u679c',
    accent: '#e35d53',
    icon: `
      <circle cx="160" cy="118" r="30" fill="#fff6ef"/>
      <circle cx="190" cy="108" r="26" fill="#ffe7d8"/>
      <path d="M170 85c8-10 18-16 31-16-6 10-15 18-29 22z" fill="#2b9a66"/>
      <rect x="168" y="72" width="6" height="20" rx="3" fill="#8b5a2b"/>
    `,
  },
  vegetable: {
    label: '\u852c\u83dc',
    accent: '#2f9e62',
    icon: `
      <path d="M160 80c-24 12-39 35-39 63 0 13 3 24 8 34 21-8 38-22 50-43 9-17 13-35 11-54-9 1-18 1-30 0z" fill="#d8f7e5"/>
      <path d="M165 88c24 12 39 35 39 63 0 13-3 24-8 34-21-8-38-22-50-43-9-17-13-35-11-54 9 1 18 1 30 0z" fill="#f0fff5"/>
      <path d="M162 92c0 48-10 79-28 102" stroke="#2f9e62" stroke-width="6" stroke-linecap="round"/>
      <path d="M162 98c0 38 12 67 30 92" stroke="#62bd86" stroke-width="6" stroke-linecap="round"/>
    `,
  },
  protein: {
    label: '\u86cb\u767d',
    accent: '#c65a36',
    icon: `
      <path d="M120 102c0-19 16-34 35-34h22c24 0 43 19 43 43v12c0 31-25 56-56 56h-8c-31 0-56-25-56-56v-21c0-6 4-10 10-10h10z" fill="#fff2ea"/>
      <circle cx="189" cy="101" r="14" fill="#ffd5c3"/>
      <path d="M110 117c15-16 31-26 54-26 27 0 45 13 56 31" stroke="#c65a36" stroke-width="8" stroke-linecap="round"/>
    `,
  },
  seafood: {
    label: '\u6d77\u9c9c',
    accent: '#0f8ca8',
    icon: `
      <path d="M102 130c19-31 50-48 88-48 22 0 45 8 64 24-8 2-15 9-18 18 4 9 10 16 18 18-19 16-42 24-64 24-38 0-69-17-88-48 5 0 10 1 14 2 4-6 4-14 0-20-4 1-9 2-14 2z" fill="#effbff"/>
      <circle cx="181" cy="118" r="5" fill="#0f8ca8"/>
      <path d="M113 111l-24-18v74l24-18" fill="#c9f0f8"/>
    `,
  },
  egg: {
    label: '\u9e21\u86cb',
    accent: '#cc8b00',
    icon: `
      <ellipse cx="160" cy="122" rx="45" ry="57" fill="#fff9ef"/>
      <circle cx="160" cy="130" r="18" fill="#ffd76f"/>
      <path d="M160 70c14 0 27 7 35 20 7 10 10 22 10 35 0 32-20 54-45 54s-45-22-45-54c0-13 3-25 10-35 8-13 21-20 35-20z" fill="none" stroke="#cc8b00" stroke-width="4" opacity=".35"/>
    `,
  },
  dairy: {
    label: '\u4e73\u54c1',
    accent: '#6d5bd0',
    icon: `
      <path d="M139 83h42l10 26v58c0 11-9 20-20 20h-22c-11 0-20-9-20-20v-58z" fill="#f4f1ff"/>
      <path d="M146 73h28l8 17h-44z" fill="#ddd5ff"/>
      <rect x="149" y="120" width="22" height="30" rx="8" fill="#c9bcff"/>
    `,
  },
  tofu: {
    label: '\u8c46\u5236\u54c1',
    accent: '#91744a',
    icon: `
      <rect x="120" y="92" width="52" height="52" rx="10" fill="#fff8e8"/>
      <rect x="150" y="112" width="52" height="52" rx="10" fill="#f6e7c5"/>
      <path d="M120 118l30-18m22 44l30-18" stroke="#d8c39a" stroke-width="4" opacity=".7"/>
    `,
  },
  staple: {
    label: '\u4e3b\u98df',
    accent: '#b36a17',
    icon: `
      <path d="M115 122c0 30 20 54 45 54s45-24 45-54z" fill="#fff4dd"/>
      <path d="M103 117h114" stroke="#b36a17" stroke-width="8" stroke-linecap="round"/>
      <path d="M136 146c12-8 24-12 37-12 17 0 29 5 39 12" stroke="#d49d55" stroke-width="6" stroke-linecap="round"/>
      <path d="M128 160c11-6 22-9 32-9 16 0 30 5 43 13" stroke="#e1b97e" stroke-width="6" stroke-linecap="round"/>
    `,
  },
  soup: {
    label: '\u6c64\u7fb9',
    accent: '#0f7d75',
    icon: `
      <path d="M113 126c0 29 21 52 47 52h8c26 0 47-23 47-52z" fill="#ebfffb"/>
      <path d="M101 121h126" stroke="#0f7d75" stroke-width="8" stroke-linecap="round"/>
      <path d="M138 90c-8 10-8 22 0 32" stroke="#70c9bf" stroke-width="5" stroke-linecap="round"/>
      <path d="M160 83c-8 12-8 24 0 36" stroke="#70c9bf" stroke-width="5" stroke-linecap="round"/>
      <path d="M182 90c-8 10-8 22 0 32" stroke="#70c9bf" stroke-width="5" stroke-linecap="round"/>
    `,
  },
  nut: {
    label: '\u575a\u679c',
    accent: '#946132',
    icon: `
      <path d="M139 94c-14 12-20 28-18 47 2 24 17 40 39 40 21 0 36-15 39-39 2-20-5-37-20-49-14 12-27 18-40 18z" fill="#fff3e6"/>
      <path d="M159 91c7 12 21 21 40 28" stroke="#946132" stroke-width="5" stroke-linecap="round" opacity=".6"/>
    `,
  },
  default: {
    label: '\u9910\u98df',
    accent: '#54677c',
    icon: `
      <circle cx="160" cy="122" r="48" fill="#f7f9fc"/>
      <path d="M160 84c23 0 42 19 42 42s-19 42-42 42-42-19-42-42 19-42 42-42z" fill="none" stroke="#a9b7c7" stroke-width="8"/>
      <circle cx="160" cy="122" r="18" fill="#dbe4ef"/>
    `,
  },
}

const CATEGORY_RULES = [
  { key: 'fruit', keywords: ['\u82f9\u679c', '\u68a8', '\u6a59', '\u6a58', '\u9999\u8549', '\u7315\u7334\u6843', '\u6843', '\u674e\u5b50', '\u8349\u8393', '\u84dd\u8393', '\u54c8\u5bc6\u74dc', '\u706b\u9f99\u679c', '\u5723\u5973\u679c', '\u67da\u5b50', 'apple', 'banana', 'berry'] },
  { key: 'nut', keywords: ['\u575a\u679c', '\u674f\u4ec1', '\u6838\u6843', '\u8170\u679c', '\u82b1\u751f', '\u699b\u5b50', '\u5f00\u5fc3\u679c', 'nut', 'almond', 'walnut', 'peanut'] },
  { key: 'dairy', keywords: ['\u8131\u8102\u5976', '\u725b\u5976', '\u9178\u5976', '\u5976\u916a', '\u5976\u6614', '\u829d\u58eb', 'milk', 'yogurt', 'cheese'] },
  { key: 'tofu', keywords: ['\u8c46\u8150', '\u8c46\u5e72', '\u9999\u5e72', '\u8150\u7af9', 'tofu'] },
  { key: 'egg', keywords: ['\u9e21\u86cb', '\u9e4c\u9e51\u86cb', '\u86cb\u7fb9', '\u8336\u53f6\u86cb', '\u6c34\u716e\u86cb', 'egg'] },
  { key: 'seafood', keywords: ['\u4e09\u6587\u9c7c', '\u9cd5\u9c7c', '\u867e', '\u87f9', '\u6d77\u5e26', '\u6d77\u9c9c', '\u9c7c', 'salmon', 'shrimp', 'fish'] },
  { key: 'protein', keywords: ['\u9e21\u80f8', '\u9e21\u817f', '\u9e21\u4e01', '\u6ed1\u9e21', '\u9e21\u4e1d', '\u9e21\u8089', '\u725b\u8089', '\u725b\u8169', '\u725b\u6392', '\u732a\u8089', '\u8089\u4e1d', '\u8089\u4e38', '\u725b\u8171', 'chicken', 'beef', 'pork', 'steak'] },
  { key: 'soup', keywords: ['\u7ca5', '\u6c64', '\u7fb9', '\u8c46\u6d46', 'soup', 'porridge'] },
  { key: 'staple', keywords: ['\u7c73\u996d', '\u9762', '\u9984\u5934', '\u5410\u53f8', '\u5305\u5b50', '\u82b1\u5377', '\u53d1\u7cd5', '\u7a9d\u5934', '\u7389\u7c73', '\u71d5\u9ea6', '\u7d2b\u85af', '\u7ea2\u85af', '\u5357\u74dc\u996d', '\u5c0f\u7c73\u996d', '\u7cd9\u7c73\u996d', 'rice', 'noodle', 'bread', 'toast', 'oat'] },
  { key: 'vegetable', keywords: ['\u897f\u5170\u82b1', '\u9ec4\u74dc', '\u83e0\u83dc', '\u83dc\u82b1', '\u83dc\u5fc3', '\u5a03\u5a03\u83dc', '\u5357\u74dc', '\u51ac\u74dc', '\u5c71\u836f', '\u6728\u8033', '\u82b9\u83dc', '\u6cb9\u9ea6\u83dc', '\u9999\u83c7', '\u897f\u846b\u82a6', '\u8c46\u89d2', '\u8377\u5170\u8c46', '\u756a\u8304', '\u897f\u7ea2\u67ff', '\u9752\u6912', 'broccoli', 'spinach', 'tomato', 'cucumber', 'pepper'] },
]

const SOUP_KEYWORDS = ['\u6c64', '\u7fb9', '\u7ca5', 'soup', 'porridge']

function escapeXml(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;')
}

function normalizeMealType(mealType) {
  const value = String(mealType || '').trim().toLowerCase()
  if (!value) return 'default'
  if (value === '\u65e9\u9910' || value === 'breakfast') return 'breakfast'
  if (value === '\u5348\u9910' || value === 'lunch') return 'lunch'
  if (value === '\u665a\u9910' || value === 'dinner') return 'dinner'
  if (value === '\u52a0\u9910' || value === 'snack') return 'snack'
  return 'default'
}

function normalizeList(value) {
  if (!value) return []
  if (Array.isArray(value)) return value.map((item) => String(item).trim()).filter(Boolean)
  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return []
    try {
      const parsed = JSON.parse(text)
      return Array.isArray(parsed) ? parsed.map((item) => String(item).trim()).filter(Boolean) : [text]
    } catch {
      return [text]
    }
  }
  return [String(value).trim()]
}

function hasKeyword(content, keywords) {
  return keywords.some((keyword) => content.includes(keyword))
}

function inferCategory(name, ingredients = []) {
  const content = `${name || ''} ${normalizeList(ingredients).join(' ')}`.toLowerCase()
  const isSoupLike = hasKeyword(content, SOUP_KEYWORDS)

  for (const rule of CATEGORY_RULES) {
    if (!hasKeyword(content, rule.keywords)) continue
    if ((rule.key === 'seafood' || rule.key === 'protein') && isSoupLike) {
      return 'soup'
    }
    return rule.key
  }

  return 'default'
}

function inferFocusLabel(targetUsers) {
  const values = normalizeList(targetUsers).join(' ').toLowerCase()
  if (!values) return '\u5747\u8861'
  if (values.includes('\u51cf\u8102') || values.includes('lean') || values.includes('fat loss')) return '\u51cf\u8102'
  if (values.includes('\u589e\u808c') || values.includes('build') || values.includes('muscle gain')) return '\u589e\u808c'
  if (values.includes('\u63a7\u7cd6') || values.includes('sugar') || values.includes('low gi')) return '\u63a7\u7cd6'
  return '\u5747\u8861'
}

function pickNameBadge(name) {
  const text = String(name || '').trim()
  return escapeXml(text ? text.slice(0, Math.min(6, text.length)) : '\u667a\u80fd\u914d\u9910')
}

function svgToDataUri(svg) {
  return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`
}

export function buildRecipePlaceholderSrc({ name, mealType, ingredients, targetUsers } = {}) {
  const meal = MEAL_META[normalizeMealType(mealType)] || MEAL_META.default
  const category = CATEGORY_META[inferCategory(name, ingredients)] || CATEGORY_META.default
  const focus = inferFocusLabel(targetUsers)
  const titleBadge = pickNameBadge(name)

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 240" preserveAspectRatio="xMidYMid meet" role="img" aria-label="${escapeXml(name || meal.label)}">
      <defs>
        <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="${meal.start}" />
          <stop offset="100%" stop-color="${meal.end}" />
        </linearGradient>
        <linearGradient id="shine" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stop-color="#ffffff" stop-opacity=".84" />
          <stop offset="100%" stop-color="#ffffff" stop-opacity=".24" />
        </linearGradient>
      </defs>
      <rect width="320" height="240" rx="28" fill="url(#bg)" />
      <circle cx="36" cy="34" r="48" fill="#ffffff" opacity=".26" />
      <circle cx="292" cy="40" r="56" fill="#ffffff" opacity=".16" />
      <circle cx="286" cy="212" r="62" fill="#ffffff" opacity=".14" />
      <rect x="18" y="18" width="84" height="36" rx="18" fill="#ffffff" opacity=".92" />
      <text x="60" y="36" text-anchor="middle" dominant-baseline="middle" font-size="18" font-family="${FONT_STACK}" font-weight="800" fill="${meal.accent}">${meal.label}</text>
      <rect x="235" y="18" width="67" height="34" rx="17" fill="${category.accent}" opacity=".94" />
      <text x="268.5" y="35" text-anchor="middle" dominant-baseline="middle" font-size="16" font-family="${FONT_STACK}" font-weight="800" fill="#ffffff">${focus}</text>
      <g transform="translate(0 2)">
        <ellipse cx="160" cy="188" rx="76" ry="18" fill="#0f172a" opacity=".09" />
        <rect x="80" y="54" width="160" height="144" rx="34" fill="url(#shine)" opacity=".7" />
        <rect x="92" y="66" width="136" height="120" rx="28" fill="#ffffff" opacity=".22" />
        <g>${category.icon}</g>
      </g>
      <rect x="20" y="182" width="144" height="34" rx="17" fill="#ffffff" opacity=".95" />
      <text x="34" y="199" dominant-baseline="middle" font-size="16" font-family="${FONT_STACK}" font-weight="800" fill="#20323a">${titleBadge}</text>
      <rect x="226" y="182" width="76" height="34" rx="17" fill="#ffffff" opacity=".92" />
      <text x="264" y="199" text-anchor="middle" dominant-baseline="middle" font-size="16" font-family="${FONT_STACK}" font-weight="800" fill="${category.accent}">${category.label}</text>
    </svg>
  `.trim()

  return svgToDataUri(svg)
}

export function getRecipeImageSrc(input = {}) {
  const src = String(input.src || input.imageUrl || input.image_url || '').trim()
  if (src) return src
  return buildRecipePlaceholderSrc(input)
}
