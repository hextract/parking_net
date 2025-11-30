import { Loader2 } from 'lucide-react'

const LoadingSpinner = ({ fullScreen = false, size = 'default' }) => {
  const sizeClasses = {
    small: 'w-4 h-4',
    default: 'w-8 h-8',
    large: 'w-12 h-12',
  }

  const spinner = (
    <div className="flex items-center justify-center">
      <Loader2 className={`${sizeClasses[size]} animate-spin text-primary-600`} />
    </div>
  )

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-gray-50">
        {spinner}
      </div>
    )
  }

  return spinner
}

export default LoadingSpinner
