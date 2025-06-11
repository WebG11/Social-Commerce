import { Inter } from 'next/font/google'
import './globals.css'
import { Header } from '@/components/header'

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'GoBy',
  description: 'A simple e-commerce application using shadcn/ui',
}


export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className={`${inter.className} flex flex-col min-h-screen`}>
          <main>
            {children}
          </main>
      </body>
    </html>
  )
}

