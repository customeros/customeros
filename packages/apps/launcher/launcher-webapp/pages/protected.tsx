import { useState, useEffect } from "react"
import { useSession } from "next-auth/react"
import Layout from "../components/layout"
import AccessDenied from "../components/access-denied"

export default function ProtectedPage() {
  const { data: session } = useSession()
  const [content, setContent] = useState()

  // Fetch content from protected route
  useEffect(() => {
    const fetchData = async () => {
      const res = await fetch("/server/registered-apps")
      const json = await res.json()
      if (json.content) {
        setContent(json.content)
      }
    }
    fetchData()
  }, [session])


  // If no session exists, display access denied message
  if (!session) {
    return (
        <Layout>
          <AccessDenied />
        </Layout>
    )
  }

  // If session exists, display content
  return (
      <Layout>
        <h1>Protected Page</h1>
        <p>
          <strong>{content ?? "\u00a0"}</strong>
        </p>
      </Layout>
  )
  // const { data: session } = useSession()
  // const [content, setContent] = useState()
  //
  // // Fetch content from protected route
  // useEffect(() => {
  //   const fetchData = async () => {
  //     const res = await fetch("http://localhost:8070/customer-os/registered-apps")
  //     const serverResponse = await res.json()
  //     if (serverResponse) {
  //       setContent(serverResponse)
  //     }
  //   }
  //   fetchData()
  // }, [session])
  //
  //
  // // If no session exists, display access denied message
  // if (!session) {
  //   return (
  //     <Layout>
  //       <AccessDenied />
  //     </Layout>
  //   )
  // }
  //
  // // If session exists, display content
  // return (
  //   <Layout>
  //     <h1>Available applications : </h1>
  //     <p>
  //       <div>
  //         {
  //           content.apps.map((appLink: { url: string; name: string}) =>
  //               <a href={appLink.url}>
  //                 <h2>{appLink.name}</h2>
  //               </a>
  //           )
  //         }
  //       </div>
  //       <div>
  //         <p>{content.content}</p>
  //       </div>
  //     </p>
  //   </Layout>
  // )
}
