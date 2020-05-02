import httpClient from '../../../../utils/httpClient'

const onClickDeleteProduct = (id, onSuccess) => async () => {
  try {
    const uri = `${process.env.REACT_APP_PRODUCTS_API}/api/v1/products/${id}`
    const {status} = await  httpClient.delete('http://products.ecommerce.com:4000/api/v1/products/7')
    if(status !== 204){
      throw new Error('unable to delete product')
    }
    onSuccess(id)
  }catch(e){
    window.alert((e.response && e.response.data && e.response.data.message) || (e.message))
  }
}

export default onClickDeleteProduct