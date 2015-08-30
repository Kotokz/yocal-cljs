(ns yocal-cljs.core-test
  (:require
   [cljs.test :refer-macros [deftest testing is]]
   [yocal-cljs.core :as core]))

(deftest fake-test
  (testing "fake description"
    (is (= 1 2))))
