(ns yocal-cljs.handlers
  (:require [re-frame.core :refer [register-handler path trim-v after dispatch debug log-ex]]
            [yocal-cljs.db :as db]
            [goog.crypt.base64 :as b64]
            [schema.core :as s]
            [ajax.core :as ajax]
            [yocal-cljs.view.login :as login]
            [yocal-cljs.view.signup :as register]))

;; -- Middleware --------------------------------------------------------------
;;

(defn check-and-throw
  "throw an exception if db doesn't match the schema."
  [a-schema data]
  (if-let [problems (s/check a-schema data)]
    (throw (js/Error. (str "schema check failed: " problems)))))

;; after an event handler has run, this middleware can check that
;; it the value in app-db still correctly matches the schema.
(def check-schema-mw (after (partial check-and-throw db/schema)))
(def dice-middleware [check-schema-mw                       ;; after ever event handler make sure the schema is still valid
                      (path :game-dices)
                      debug
                      trim-v                                ;; remove event id from event vec
                      ])
(def auth-middleware [check-schema-mw                       ;; after ever event handler make sure the schema is still valid
                      (path :user)
                      debug
                      trim-v                                ;; remove event id from event vec
                      ])

;; -- Helpers -----------------------------------------------------------------
(defn auth-header [crets]
  (let [{:keys [username password]} crets]
    (pr username password)
    (str "Basic " (b64/encodeString (str username ":" password)))))

;; -- Handlers ----------------------------------------------------------------------

(register-handler
  :initialize-db
  check-schema-mw
  (fn [_ _]
    db/default-db))

(register-handler
  :set-active-panel
  debug
  (fn [db [_ active-panel]]
    (assoc db :active-panel active-panel)))

(register-handler
  :set-roll-count
  [(path :game-rolls)]
  (fn [game-rolls _]
    (inc game-rolls)))

(register-handler
  :add-dice
  dice-middleware
  (fn [game-dices [id]]
    (let [val (inc (rand-int 6))]
      (assoc game-dices id {:id id :val val :isHeld false}))))

(register-handler
  :roll-dice
  dice-middleware
  (fn [game-dices [id]]
    (let [val (inc (rand-int 6))]
      (assoc-in game-dices [id :val] val))))

(register-handler
  :hold-dice
  dice-middleware
  (fn [game-dices [id]]
    (update-in game-dices [id :isHeld] not)))

(register-handler
  :set-score-string
  [(path :game-score-string) debug trim-v]
  (fn [game-score-string [text]]
    text))

;;-- Login hanlder ----------------------------------------------------------------------
(register-handler
  :login
  auth-middleware
  (fn [user [crets]]
    (swap! login/login-form-data update-in [:is-loading] not)
    (ajax/POST "http://127.0.0.1:3001/user/token" {:headers         {;"Authorization" (auth-header crets)
                                                                "X-Requested-With" "XMLHttpRequest"}
                                              :params          crets
                                              :format          :json
                                              :response-format :json
                                              :keywords?       true
                                              :prefix          true
                                              :handler         #(dispatch [:login-response %1])
                                              :error-handler   #(dispatch [:login-failed %1])})
    user))

(register-handler
  :login-response
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (swap! login/login-form-data update-in [:is-loading] not)
    (assoc user :jwt (:jwt response))))

(register-handler
  :login-failed
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (let [errs (get-in response [:response :errors])]
      (swap! login/login-form-data update-in [:is-loading] not)
      ;(swap! login/login-form-data assoc-in [:errors :other] (get-in response [:response :msg]))
      (swap! login/login-form-data assoc :errors errs)
      (swap! login/login-form-data assoc-in [:errors :other] (get-in response [:response :msg])))
    user))

;;-- Register hanlder ----------------------------------------------------------------------

(register-handler
  :sign-up
  auth-middleware
  (fn [user [crets]]
    (swap! register/signup-form-data update-in [:is-loading] not)
    (ajax/POST "http://127.0.0.1:3001/user/register" {:headers         {;"Authorization" (auth-header crets)
                                                                        "X-Requested-With" "XMLHttpRequest"}
                                                      :params          crets
                                                      :format          :json
                                                      :response-format :json
                                                      :keywords?       true
                                                      :prefix          true
                                                      :handler         #(dispatch [:sign-up-response %1])
                                                      :error-handler   #(dispatch [:sign-up-failed %1])})
    user))

(register-handler
  :sign-up-response
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (swap! register/signup-form-data update-in [:is-loading] not)
    (pr response)
    user))

(register-handler
  :sign-up-failed
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (let [errs (get-in response [:response :errors])]
      (swap! register/signup-form-data update-in [:is-loading] not)
      (swap! register/signup-form-data assoc :errors errs)
      (swap! register/signup-form-data assoc-in [:errors :other] (get-in response [:response :msg])))
    user))

;;-- API hanlder ----------------------------------------------------------------------

(register-handler
  :get-balance
  auth-middleware
  (fn [user [jwt]]
    (ajax/POST "http://127.0.0.1:3001/api/balance" {:headers         {"Authorization" (str "Bearer " jwt)}
                                                     :response-format :json
                                                     :keywords?       true
                                                     :prefix          true
                                                     :handler         #(dispatch [:get-balance-ok %1])
                                                     :error-handler   #(dispatch [:get-balance-failed %1])})
    user))

(register-handler
  :get-balance-ok
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (pr (:balance response))
    user))

(register-handler
  :get-balance-failed
  (fn
    ;; store info for the specific phone-id in the db
    [db [_ response]]
    (pr response)
    db))
