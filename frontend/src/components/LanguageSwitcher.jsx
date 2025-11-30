import { useTranslation } from 'react-i18next'
import { Languages } from 'lucide-react'

const LanguageSwitcher = () => {
  const { i18n } = useTranslation()

  const toggleLanguage = async () => {
    const newLang = i18n.language === 'en' ? 'ru' : 'en'
    await i18n.changeLanguage(newLang)
    localStorage.setItem('language', newLang)
  }

  return (
    <button
      onClick={toggleLanguage}
      className="flex items-center space-x-2 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors"
      title={i18n.language === 'en' ? 'Switch to Russian' : 'Переключить на английский'}
    >
      <Languages className="w-4 h-4" />
      <span className="text-sm font-medium uppercase">{i18n.language}</span>
    </button>
  )
}

export default LanguageSwitcher
