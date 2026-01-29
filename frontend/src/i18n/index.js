import { createI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import 'dayjs/locale/zh-cn'
import 'dayjs/locale/en'
import en from './locales/en.json'
import zh from './locales/zh.json'

dayjs.extend(localizedFormat)


const initialLocale = localStorage.getItem('lang') || 'en'

// Set initial dayjs locale
dayjs.locale(initialLocale === 'zh' ? 'zh-cn' : 'en')

const i18n = createI18n({
    legacy: false,
    globalInjection: true,
    locale: initialLocale,
    fallbackLocale: 'en',
    messages: {
        en,
        zh
    }
})

export default i18n