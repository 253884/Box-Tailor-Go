# Box-Tailor-Go

GUI application made to simplify the process of creating tailored boxes for your products.
The app is written in `Go` using `Sciter`. 
App settings are stored using a simple `SQLite` database.

## Usage:

Main window:

![box_tailor_1](https://user-images.githubusercontent.com/73070465/115378495-e60fd980-a1d0-11eb-9b01-8eb9143fa68d.png)

Adding box from file (`*.plt` files supported only. More filetypes are to be expected in [a follow-up project](https://github.com/happyRip/Box-Tailor)):

![box_tailor_2](https://user-images.githubusercontent.com/73070465/115378850-49017080-a1d1-11eb-8b23-35ee4fb833e3.png)

Select type of box (either with integrated lid or a standard flap box):

![box_tailor_3](https://user-images.githubusercontent.com/73070465/115378891-528ad880-a1d1-11eb-8935-947ebebc7337.png)

Adding box by product dimensions (if your file is incompatible):

![box_tailor_4](https://user-images.githubusercontent.com/73070465/115379195-97167400-a1d1-11eb-9e3f-9d156f4040a1.png)

## Output based on parameters:

Flap boxes (low profile product, mailer would offer a better fit):

![box_tailor_plt_1](https://user-images.githubusercontent.com/73070465/115379216-9d0c5500-a1d1-11eb-805b-7c170034af36.png)

Message boxes:

![box_tailor_plt_2](https://user-images.githubusercontent.com/73070465/115379219-9d0c5500-a1d1-11eb-8f38-12a11c653c34.png)

## Note: 

Project **abandoned**. Working with `Sciter` was not as pleasant as using pure Go and Web tools. Nonetheless the tool will be accessible on my website in the future utilizing [the follow-up project](https://github.com/happyRip/Box-Tailor).

