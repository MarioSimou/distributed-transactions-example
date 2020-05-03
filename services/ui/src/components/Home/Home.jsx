import React from 'react'
import {
  Paper,
  Typography,
  Toolbar,
  IconButton,
} from '@material-ui/core'
import {makeStyles} from '@material-ui/core'
import AddIcon from '@material-ui/icons/Add'
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';
import ProductCard from './components/ProductCard'
import fetchProducts from './utils/fetchProducts.js'
import {CancelToken} from 'axios'
import AddProductDialog from './components/AddProductDialog'
import BuyNow from './components/BuyNow'
import {useUserProfile} from '../../utils/hooks'
import onClickBuyNow from './utils/onClickBuyNow'

const Home = () => {
  const classes = useStyles()
  const {userProfile} = useUserProfile()
  const [showDialog, setShowDialog] = React.useState(false)
  const [products, setProducts] = React.useState([])
  const [userProducts, setUserProducts] = React.useState([])
  const onCloseDialog = () => setShowDialog(false)
  const onClickAddProduct = () => setShowDialog(true)
  const onDeleteProduct = id => products.filter(product => product.id !== id) 
  const onAddToCart = product => () => setUserProducts([...userProducts, product])
  const resetCart = () => setUserProducts([])

  React.useEffect(() => {
    const source = CancelToken.source()
    fetchProducts(setProducts, source)
    return () => source.cancel()
  },[])

  return (
    <Paper className={classes.root}>
      <Toolbar variant="dense" className={classes.toolbar}>
        <IconButton className={classes.iconButton} onClick={onClickAddProduct}>
          <AddIcon/>
        </IconButton>
        <IconButton className={classes.iconButton}>
          <Typography component="span" className={classes.shoppingCartCounter}>{userProducts.length}</Typography>
          <ShoppingCartIcon fontSize="small"/>
        </IconButton>
      </Toolbar>
      <Typography className={classes.main} component="div">
        <AddProductDialog open={showDialog} onCloseDialog={onCloseDialog} />
        <Typography component="div" className={classes.productsWrapper}>
          {products.length && products.map(product => {
            return <ProductCard key={product.id} {...product} onDeleteProduct={onDeleteProduct} onAddToCart={onAddToCart(product)}/>
          })}
        </Typography>
        {userProducts.length > 0 && <BuyNow className={classes.buyNow} onClickBuyNow={onClickBuyNow(userProducts, userProfile, resetCart)}/>}
      </Typography>
    </Paper>
  )
}

const useStyles = makeStyles(theme => ({
  root: {
    minHeight: 'calc( 100vh - 48px )',
    position: 'relative',
  },
  toolbar: {
    width: '100%',
    backgroundColor: '#e8e8e8',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
  },
  main: {
    position: 'relative',
    height: 'calc( 100vh - 98px )',
    padding: theme.spacing(4),
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(2),
    }
  },
  productsWrapper: {
    display: 'grid',
    gridGap: theme.spacing(2),
    gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
  },
  iconButton: {
    position: 'relative',
  },
  shoppingCartCounter: {
    position: 'absolute',
    fontWeight: 600,
    top: -theme.spacing(2),
    padding: theme.spacing(1.5),
    overflow: 'visible',
  },
  buyNow: {
    position: 'absolute',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
    '& button:last-child': {
      marginLeft: theme.spacing(1),
    } 
  }
}))

export default Home