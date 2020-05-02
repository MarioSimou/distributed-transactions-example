import httpClient from '../../../utils/httpClient'

const fetchProducts = async (setProducts,source) => {
  try {
    const uri = new URL('/api/v1/products', process.env.REACT_APP_PRODUCTS_API)
    const {data} = await httpClient.get(uri.toString(), {cancelToken: source.token})
    if(!data.success){
      throw new Error(data.message)
    }
    setProducts(data.data)
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message))
  }
}

export default fetchProducts