package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

var jwtkey = []byte("appscode")
var tokenAuth *jwtauth.JWTAuth
var tokenString string
var token jwt.Token

type Products struct {
	ProductId     int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	ProductBrand  *Brands `json:"product_brand"`
	ProductSize   string  `json:"product_size"`
	ProductPrice  int     `json:"product_price"`
	ProductStatus bool    `json:"product_status"`
}

type Brands struct {
	BrandId      int    `json:"brand_id"`
	BrandName    string `json:"brand_name"`
	BrandProduct []int  `json:"brand_product"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CredsDB map[string]string
type StuffDB map[int]Products
type BrandDB map[int]*Brands

var productlist StuffDB
var brandlist BrandDB
var credslist CredsDB

func Slice_rm(x []int, val int) []int {
	for i := 0; i < len(x); i++ {
		if x[i] == val {
			x = append(x[:i], x[i+1:]...)
			break
		}
	}
	return x
}

func init() {
	productlist = make(StuffDB)
	brandlist = make(BrandDB)
	credslist = make(CredsDB)
	tokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
}

func productGen() {
	productlist[1] = Products{
		ProductId:    1,
		ProductName:  "Shirt",
		ProductBrand: brandlist[1],
		ProductSize:  "M",
		ProductPrice: 240,
	}

	productlist[2] = Products{
		ProductId:    2,
		ProductName:  "Pant",
		ProductBrand: brandlist[2],
		ProductSize:  "M",
		ProductPrice: 280,
	}
}

func brandGen() {
	brandlist[1] = &Brands{
		BrandId:      1,
		BrandName:    "Gucci",
		BrandProduct: []int{1},
	}
	brandlist[2] = &Brands{
		BrandId:      2,
		BrandName:    "anyBrand",
		BrandProduct: []int{2},
	}
}

func GenCreds() {
	creds := credentials{
		"Rakibul Hossain",
		"123456",
	}
	credslist[creds.Username] = creds.Password

}

func GenData() {
	brandGen()
	productGen()
	GenCreds()
}

func GetBrandId() int {
	avail := 1
	for value, _ := range brandlist {
		if value > avail {
			return avail
		} else {
			avail++
		}
	}
	return avail
}

func GetProductId() int {
	avail := 1
	for value, _ := range productlist {
		if value > avail {
			return avail
		} else {
			avail++
		}
	}
	return avail
}

func ShowAllBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(brandlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ShowAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(productlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productid := chi.URLParam(r, "productsid")
	id, _ := strconv.Atoi(productid)

	if _, ok := productlist[id]; !ok {
		w.WriteHeader(404)
		return
	}
	err := json.NewEncoder(w).Encode(productlist[id])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brandid := chi.URLParam(r, "brandsid")
	id, _ := strconv.Atoi(brandid)

	if _, ok := brandlist[id]; !ok {
		w.WriteHeader(404)
		return
	}
	err := json.NewEncoder(w).Encode(brandlist[id])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	brandid := chi.URLParam(r, "brandsid")
	id, _ := strconv.Atoi(brandid)
	if _, ok := brandlist[id]; !ok {
		w.WriteHeader(404)
		return
	}
	if len(brandlist[id].BrandProduct) > 0 {
		w.Write([]byte("You have products of this brand"))
		return
	}
	delete(brandlist, id)
	err := json.NewEncoder(w).Encode(brandlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productid := chi.URLParam(r, "productsid")
	id, _ := strconv.Atoi(productid)
	if _, ok := productlist[id]; !ok {
		w.WriteHeader(404)
		return
	}
	temp := productlist[id].ProductBrand
	temp.BrandProduct = Slice_rm(temp.BrandProduct, id)
	delete(productlist, id)
	err := json.NewEncoder(w).Encode(productlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// TODO: add some resource
// TODO: Documentation

func AddBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBrand Brands
	err := json.NewDecoder(r.Body).Decode(&newBrand)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newBrand.BrandId = GetBrandId()
	brandlist[newBrand.BrandId] = &newBrand
	err = json.NewEncoder(w).Encode(brandlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func AddProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newProduct Products
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newProduct.ProductId = GetProductId()
	temp := newProduct.ProductBrand
	ok := false
	for key, value := range brandlist {
		if value == temp {
			brandlist[key].BrandProduct = append(brandlist[key].BrandProduct, newProduct.ProductId)
			ok = true
			newProduct.ProductBrand = brandlist[key]
			break
		}
	}
	if !ok {
		newBrand := &Brands{
			BrandId:      1,
			BrandName:    newProduct.ProductBrand.BrandName,
			BrandProduct: []int{newProduct.ProductId},
		}
		newBrand.BrandId = GetBrandId()
		fmt.Println(newBrand.BrandId)
		brandlist[newBrand.BrandId] = newBrand
		newProduct.ProductBrand = newBrand
	}
	productlist[newProduct.ProductId] = newProduct
	err = json.NewEncoder(w).Encode(productlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func UpdateProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productid := chi.URLParam(r, "productsid")

	id, _ := strconv.Atoi(productid)
	if _, ok := productlist[id]; !ok {
		w.WriteHeader(404)
		return
	}
	newProduct := productlist[id]
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	productlist[id] = newProduct
	err = json.NewEncoder(w).Encode(productlist)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var creds credentials

	err := json.NewDecoder(r.Body).Decode(&creds)

	fmt.Println(creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	correctPassword, ok := credslist[creds.Username]

	if !ok || creds.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//tokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)

	expiretime := time.Now().Add(15 * time.Minute)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"aud": "Rakibul Hossain",
		"exp": expiretime.Unix(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: expiretime,
	})
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

func main() {

	GenData()
	r := chi.NewRouter()
	fmt.Println("started")
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/login", Login)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		//	r.Use(middleware.BasicAuth("user", credslist))

		r.Route("/products", func(r chi.Router) {
			r.Get("/", ShowAllProducts)
			r.Get("/{productsid}", GetProducts)
			r.Post("/add", AddProducts)
			r.Post("/delete/{productsid}", DeleteProducts)
			r.Post("/update/{productsid}", UpdateProducts)
		})

		r.Route("/brands", func(r chi.Router) {
			r.Get("/", ShowAllBrands)
			r.Get("/{brandsid}", GetBrands)
			r.Post("/add", AddBrands)
			r.Post("/delete/{brandsid}", DeleteBrands)
		})
		r.Post("/logout", Logout)

	})
	http.ListenAndServe(":8081", r)
}
