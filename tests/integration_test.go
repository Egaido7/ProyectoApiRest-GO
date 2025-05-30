package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"parte3/api"
	"parte3/internal/sale"
	"parte3/internal/user"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

// validación de main.go para que tome el regex
func regexpValidationTest(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// La misma regex que en main.go
	regex := regexp.MustCompile(`^[a-zA-Z]+(?: [a-zA-Z]+)*$`)
	return regex.MatchString(value)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Registrar la validación personalizada para el motor de prueba
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("regexp", regexpValidationTest)
	}

	api.InitRoutes(router) // inicializar tus servicios y rutas
	return router
}

// generarString crea una cadena aleatoria de letras de una longitud dada.
func generarString(length int) string {
	rand.Seed(time.Now().UnixNano())
	letras := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	var b strings.Builder //construir cadena concatenando las letras
	for i := 0; i < length; i++ {
		b.WriteRune(letras[rand.Intn(len(letras))])
	}
	return b.String()
}

// crearUsuarioforTest es una función auxiliar para crear un usuario vía API y devolver su ID.
// Asegura que el nombre de usuario cumpla con la regexp `^[a-zA-Z]+(?: [a-zA-Z]+)*$`.
func crearUsuarioforTest(t *testing.T, router *gin.Engine) string {
	userName := "TestUser" + generarString(4) // ej., TestUserXyzAbc
	userPayload := gin.H{
		"name":    userName,
		"address": "25 de mayo 1234",
	}
	userBody, err := json.Marshal(userPayload) // lo hago para serializar la info de user a json
	require.NoError(t, err)

	reqPostUser, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(userBody))
	reqPostUser.Header.Set("Content-Type", "application/json")
	recorderUser := httptest.NewRecorder()      //emulador de respuesta HTTP
	router.ServeHTTP(recorderUser, reqPostUser) // permite al enrutador (router) manejar solicitudes HTTP.

	if recorderUser.Code != http.StatusCreated {
		t.Logf("Cuerpo de respuesta al crear usuario: %s", recorderUser.Body.String())
	}
	require.Equal(t, http.StatusCreated, recorderUser.Code, "Fallo al crear usuario para la prueba")

	var createdUser user.User
	err = json.Unmarshal(recorderUser.Body.Bytes(), &createdUser) // Deserializo la respuesta JSON a un objeto User
	// Verifico que la deserialización no falle y que el ID del usuario creado no esté vacío
	require.NoError(t, err)
	require.NotEmpty(t, createdUser.ID, "El ID del usuario creado está vacío") //aserciones por si falla detiene la ejecucion
	return createdUser.ID
}

// HappyPath_PostPatchGet prueba la secuencia completa POST -> PATCH -> GET para las ventas.
func HappyPath_PostPatchGet(t *testing.T) {
	router := setupRouter()
	testUserID := crearUsuarioforTest(t, router)

	// 1. POST /sales (Crear Venta)
	saleAmount := 300.50
	createSalePayload := sale.CreateSaleRequest{
		UserID: testUserID,
		Amount: saleAmount, // El monto debe ser > 0 según el binding de sale.CreateSaleRequest
	}
	saleBody, err := json.Marshal(createSalePayload)
	require.NoError(t, err)

	reqPostSale, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewBuffer(saleBody))
	reqPostSale.Header.Set("Content-Type", "application/json")
	rrPostSale := httptest.NewRecorder()
	router.ServeHTTP(rrPostSale, reqPostSale)

	if rrPostSale.Code != http.StatusCreated {
		t.Logf("POST /sales estado de respuesta: %d, cuerpo: %s", rrPostSale.Code, rrPostSale.Body.String())
	}
	require.Equal(t, http.StatusCreated, rrPostSale.Code, "Fallo al crear la venta")

	var createdSale sale.Sale
	err = json.Unmarshal(rrPostSale.Body.Bytes(), &createdSale)
	require.NoError(t, err)
	require.NotEmpty(t, createdSale.ID, "El ID de la venta creada está vacío")
	require.Equal(t, testUserID, createdSale.UserID)
	require.Equal(t, saleAmount, createdSale.Amount)
	t.Logf("Venta creada con ID: %s, Estado Inicial: %s", createdSale.ID, createdSale.Status)

	// 2. PATCH /sales/:id (Actualizar Estado de la Venta)
	// El "camino feliz" para PATCH requiere que la venta esté en estado 'pending'.
	// se actualiza a "approved".
	StatusUpdate := "approved" // Estado objetivo válido para la actualización

	updateSalePayload := sale.UpdateSale{
		Status: StatusUpdate,
	}
	updateBody, err := json.Marshal(updateSalePayload)
	require.NoError(t, err)

	patchURL := fmt.Sprintf("/sales/%s", createdSale.ID)
	reqPatchSale, _ := http.NewRequest(http.MethodPatch, patchURL, bytes.NewBuffer(updateBody))
	reqPatchSale.Header.Set("Content-Type", "application/json")
	rrPatchSale := httptest.NewRecorder()
	router.ServeHTTP(rrPatchSale, reqPatchSale)

	var updatedSale sale.Sale //usar en el paso GET

	if createdSale.Status != "pending" {
		// Si la venta no estaba en 'pending', PATCH debería fallar con StatusConflict.
		// El manejador devuelve http.StatusConflict para ErrSaleMustBePending.
		t.Logf("El estado inicial de la venta era '%s'. Se espera que PATCH a '%s' falle.", createdSale.Status, StatusUpdate)
		require.Equal(t, http.StatusConflict, rrPatchSale.Code, "PATCH /sales no devolvió StatusConflict para una venta no pendiente. Cuerpo: "+rrPatchSale.Body.String())

		var errResp gin.H
		err = json.Unmarshal(rrPatchSale.Body.Bytes(), &errResp)
		require.NoError(t, err)
		require.Contains(t, errResp["error"], sale.ErrSaleMustBePending.Error(), "El mensaje de error para la actualización no pendiente es incorrecto.")

		t.Logf("PATCH falló como se esperaba porque el estado inicial era '%s'. El camino feliz completo POST->PATCH Exitoso->GET no puede completarse en esta ejecución.", createdSale.Status)
		return // Finaliza la prueba aquí ya que el resto de la secuencia del "camino feliz" (PATCH exitoso & GET) no puede continuar.
	}

	// Si createdSale.Status era 'pending', por lo que Actualizar(patch) debería tener éxito.
	require.Equal(t, http.StatusOK, rrPatchSale.Code, "Fallo al actualizar el estado de la venta para una venta pendiente. Cuerpo: "+rrPatchSale.Body.String())
	err = json.Unmarshal(rrPatchSale.Body.Bytes(), &updatedSale)
	require.NoError(t, err)
	require.Equal(t, StatusUpdate, updatedSale.Status, "El estado de la venta después de PATCH no es el esperado.")
	require.Equal(t, createdSale.Version+1, updatedSale.Version, "La versión de la venta no se incrementó después de PATCH.")
	t.Logf("Estado de la venta actualizado exitosamente (PATCH) a: %s", updatedSale.Status)

	// 3. GET /sales/:id/:status (Verificar Venta)
	// se ejecuta si PATCH fue exitoso (es decir, la venta estaba inicialmente en 'pending').
	// Obtendremos la venta usando su ID y el nuevo estado "approved".
	getURL := fmt.Sprintf("/sales/%s/%s", testUserID, StatusUpdate)
	reqGetSale, _ := http.NewRequest(http.MethodGet, getURL, nil)
	rrGetSale := httptest.NewRecorder()
	router.ServeHTTP(rrGetSale, reqGetSale)

	if rrGetSale.Code != http.StatusOK {
		t.Logf("GET /sales/%s/%s estado de respuesta: %d, cuerpo: %s", createdSale.ID, StatusUpdate, rrGetSale.Code, rrGetSale.Body.String())
	}
	require.Equal(t, http.StatusOK, rrGetSale.Code, "Fallo al obtener (GET) la venta por ID y estado después de la actualización.")

	// La respuesta para GET /sales/:id/:status es gin.H{"metadata": ..., "results": ...}
	var getResponse struct {
		Metadata *sale.Metadata `json:"metadata"`
		Results  []*sale.Sale   `json:"results"`
	}
	err = json.Unmarshal(rrGetSale.Body.Bytes(), &getResponse)
	require.NoError(t, err)

	require.NotNil(t, getResponse.Results, "El campo Results en la respuesta GET es nulo.")
	require.Len(t, getResponse.Results, 1, "Se esperaba una venta en el array de resultados.")

	retrievedSale := getResponse.Results[0]
	require.Equal(t, createdSale.ID, retrievedSale.ID, "El ID de la venta recuperada no coincide.")
	require.Equal(t, StatusUpdate, retrievedSale.Status, "El estado de la venta recuperada no es 'approved'.")
	require.Equal(t, updatedSale.Version, retrievedSale.Version, "La versión de la venta recuperada no coincide con la versión actualizada.")

	// Verificar metadata
	require.NotNil(t, getResponse.Metadata, "El campo Metadata en la respuesta GET es nulo.")
	require.Equal(t, 1, getResponse.Metadata.Quantity, "La cantidad en Metadata es incorrecta.")
	require.Equal(t, saleAmount, getResponse.Metadata.TotalAmount, "El monto total en Metadata es incorrecto.")
	if StatusUpdate == "approved" {
		require.Equal(t, 1, getResponse.Metadata.Approved, "El conteo de aprobadas en Metadata es incorrecto.")
		require.Equal(t, 0, getResponse.Metadata.Pending, "El conteo de pendientes en Metadata debería ser 0.")
		require.Equal(t, 0, getResponse.Metadata.Rejected, "El conteo de rechazadas en Metadata debería ser 0.")
	}
	t.Logf("Venta recuperada exitosamente con estado: %s y verificada.", retrievedSale.Status)
}
