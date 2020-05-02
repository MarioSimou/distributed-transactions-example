import React from 'react'
import {
  Paper,
  Typography,
  Toolbar,
  IconButton,
} from '@material-ui/core'
import {makeStyles} from '@material-ui/core'
import AddIcon from '@material-ui/icons/Add'
import ProductCard from './components/ProductCard'
import fetchProducts from './utils/fetchProducts.js'
import {CancelToken} from 'axios'

import AddProductDialog from './components/AddProductDialog'
const Home = () => {
  const classes = useStyles()
  const [showDialog, setShowDialog] = React.useState(false)
  const [products, setProducts] = React.useState([])
  const onCloseDialog = () => setShowDialog(false)
  const onClickAddProduct = () => setShowDialog(true)

  React.useEffect(() => {
    const source = CancelToken.source()
    fetchProducts(setProducts, source)
    return () => source.cancel()
  },[])

  console.log(products)
  return (
    <Paper className={classes.root}>
      <Toolbar variant="dense" className={classes.toolbar}>
        <IconButton size="small" onClick={onClickAddProduct}>
          <AddIcon/>
        </IconButton>
      </Toolbar>
      <Typography className={classes.main} component="div">
        <AddProductDialog open={showDialog} onCloseDialog={onCloseDialog} />
        <Typography component="div">
          {products.length && products.map(product => {
            return <ProductCard key={product.id} title={product.name} content={product.description} image={product.image}/>
          })}
        </Typography>
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
    height: '100%',
    padding: theme.spacing(4),
    [theme.breakpoints.down('sm')]: {
      padding: theme.spacing(2),
    }
  }
}))

export default Home