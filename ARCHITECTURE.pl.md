Obok głównych celów projektu podanych w [README](README.md#project-goals) 
Pixiq ma też cele poboczne, głównie te natury architektonicznej. 

Użytkownicy od tego typu bibliotek oczekują:

+ poprawnego działania
+ stabilnego API
+ elastyczności w podmienianiu wybranych elementów na inne, dodawaniu nowych

Twórcy Pixiq aby zrealizować te cele mogą:

+ tworzyć testy automatyczne (przy czym jednostkowych testów powinno być 95%)
+ starannie projektować API 
  + tworzyć tzw. proof-of-concept
  + stosować TDD
  + stosować wzorce architektoniczne, dobre praktyki a przede wszystkim
    **architekturą heksagonalną**
  + podglądać konkurencyjne rozwiązania
  + dyskutować rozwiązania z innymi
+ tworzyć małe, możliwie **niezależne** pakiety - zobacz pakiety:
  + [pixiq]()
  + [pixiq.keyboard](keyboard)
+ w większości wypadków nowe funkcje powinny się znaleźć w nowych pakietach 
(no chyba, żę czegoś brakowało od początku ;))
+ Pixiq powinien być bardziej **biblioteką** niż frameworkiem. Oznacza to, że to
  programista gier decyduje jak poskładać wszystko do kupy, nie odwrotnie. 
  Czasem wygodne jest wykorzystanie gotowych funkcji setupujących, jednak powinny 
  to być tylko dodatkowo dostępne funkcje
+ Jak tylko Pixiq osiągnie wersję 1.0.0 nie możliwa będzie zmiana API - zarówno
  syntaktyczna jak i sementaczna. Każda zmiana będzie wymagała stworzenia
  osobnego pakietu v2, v3 itp. 
+ Dlatego też z czasem zaistnieje potrzeba podziału projektu na mniejsze - tak, 
  żeby każdy moduł miał swoje własne wersjonowanie. Dzięki temu będzie można
  wprowadzać niekompatybilne zmiany do niestabilnych modułów (np. devtools),
  jednak podstawowe moduły nie będą się już zmieniać.

### Uzasadnienie dotychczasowych decyzji projektowych

Dlaczego obsługa klawiatury nie jest cześćią pakietu [pixiq]()?

> Ponieważ gra może wcale nie wykorzystać klawiatury lub działać na urządzeniach 
bez klawiatury.

Dlaczego nie ma abstrakcji na otwieranie okien?

> Ponieważ bardzo ciężko jest zaprojektować taką abstrakcję ze względu na ogromną
różnorodność platform (PC, Mac, urządzenia mobilne a kto wie może nawet konsole ;)).
Jednak jesteśmy otwarci na propozycje. Być może można stworzyć abstrakcję
tylko dla PC czyli Win, Mac i Linux.

Pakiet [opengl](opengl) wykorzystuje [GLFW](https://www.glfw.org/), a jednak nie udostępnia wszystkich
funkcji tej biblioteki np. możliwości zmiany wyglądu kursora, ustawiania
pół-przezroczystości okna itp. Czemu?

> Ponieważ nie mieliśmy jeszcze czasu tego zrobić :) Jak tylko czegoś Ci brakuje
zachęcamy do zgłoszenia [Issue](https://github.com/jacekolszak/pixiq/issues) 
lub nawet [Pull Request](https://github.com/jacekolszak/pixiq/pulls) :) Struktura `opengl.Windows` nie 
musi implementować żadnych abstrakcji, więc można ją dowolnie rozbudowywać. 
