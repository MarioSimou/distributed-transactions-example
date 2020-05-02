import httpClient from '../../../utils/httpClient.js'

export default async ({setUserProfile, source}) => {
  try {
    const uri = new URL(`/api/v1/signin`, process.env.REACT_APP_CUSTOMERS_API)
    const {data} = await httpClient.get(uri.toString(), {cancelToken: source.token})
    setUserProfile(data.data)
  }catch(e){
    console.warn('No user profile loaded')
  }
}