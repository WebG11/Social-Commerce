import { Header } from '@/components/header'

export default function RootLayout({ children }) {
  return (
    <>
        <Header />
        <div className="container mx-auto py-6 px-4 flex-grow">
          {children}
        </div>
    </>
  )
}