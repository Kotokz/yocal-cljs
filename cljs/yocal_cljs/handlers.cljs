(ns yocal-cljs.handlers
    (:require [re-frame.core :refer [register-handler path trim-v after dispatch debug log-ex]]
              [yocal-cljs.db :as db]
              [goog.crypt.base64 :as b64]
              [schema.core   :as s]
              [ajax.core :as ajax]
              [yocal-cljs.view.login :as login]))

;; -- Middleware --------------------------------------------------------------
;;

(defn check-and-throw
      "throw an exception if db doesn't match the schema."
      [a-schema data]
      (if-let [problems  (s/check a-schema data)]
         (throw (js/Error. (str "schema check failed: " problems)))))

;; after an event handler has run, this middleware can check that
;; it the value in app-db still correctly matches the schema.
(def check-schema-mw (after (partial check-and-throw db/schema)))
(def dice-middleware [check-schema-mw ;; after ever event handler make sure the schema is still valid
                      (path :game-dices)
                      debug
                      trim-v          ;; remove event id from event vec
                      ])
(def auth-middleware [check-schema-mw ;; after ever event handler make sure the schema is still valid
                      (path :user)
                      debug
                      trim-v          ;; remove event id from event vec
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
 (fn  [_ _]
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

(register-handler
  :login
  auth-middleware
  (fn [user [crets]]
    (ajax/POST "http://127.0.0.1:3001/user/token" {:headers        {;"Authorization" (auth-header crets)
                                                                     "X-Requested-With" "XMLHttpRequest"}
                                                   :params crets
                                                   :format           :json
                                                   :response-format  :json
                                                   :keywords?        true
                                                   :prefix           true
                                                   :handler          #(dispatch [:login-response %1])
                                                   :error-handler    #(dispatch [:login-failed %1])})
    user))

(register-handler
  :login-response
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (assoc user :jwt (:jwt response))))

(register-handler
  :login-failed
  (fn
    ;; store info for the specific phone-id in the db
    [db [_ response]]
    (swap! login/login-form-data assoc-in [:errors :other] (get-in response [:response :msg]))
    db))

(register-handler
  :get-balance
  auth-middleware
  (fn [user [jwt]]
    (ajax/POST "http://127.0.0.1:3001/user/balance" {:headers        {"Authorization" (str "Bearer " jwt)}
                                                   :response-format  :json
                                                   :keywords?        true
                                                   :prefix           true
                                                   :handler          #(dispatch [:get-balance-ok %1])
                                                   :error-handler    #(dispatch [:get-balance-failed %1])})
    user))

(register-handler
  :get-balance-ok
  auth-middleware
  (fn
    ;; store info for the specific phone-id in the db
    [user [response]]
    (let [rsp (js/jwt_decode (:jwt response))]
    (pr  (.. rsp -balance))
;    (pr (b64/decodeString (second (clojure.string/split (:jwt response) #"\." ) )))
    user)))

(register-handler
  :get-balance-failed
  (fn
    ;; store info for the specific phone-id in the db
    [db [_ response]]
    (pr response)
    db))
